package socket

import (
	"errors"
	"fmt"
	"github.com/FocusCompany/backend-go/api"
	"github.com/FocusCompany/backend-go/database"
	"github.com/FocusCompany/backend-go/models"
	"github.com/FocusCompany/backend-go/proto"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	zmq "github.com/pebbe/zmq4"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
	"time"
)

func getValidEvent(received string) (*Focus.Envelope, uuid.UUID, []models.Group, error) {
	envelope := &Focus.Envelope{}
	err := proto.Unmarshal([]byte(received), envelope)
	if err != nil {
		return nil, uuid.UUID{}, nil, errors.New("proto.Unmarshal error " + err.Error())
	}

	// Check JWT from envelope
	claims, err := api.ValidateJwt(envelope.Jwt)
	if err != nil {
		return nil, uuid.UUID{}, nil, errors.New("Invalid JWT " + err.Error())
	}

	userId := uuid.FromStringOrNil(claims["uuid"].(string))

	var groups []models.Group

	if claims["groups"] != nil {
		groupsClaim := claims["groups"].([]interface{})
		for _, group := range groupsClaim {
			groupId := group.(map[string]interface{})["id_collections"].(float64) // Type assertion of hell
			groupName := group.(map[string]interface{})["collections_name"].(string) // Type assertion of hell
			groups = append(groups, models.Group{
				ID: int32(groupId),
				Name:    groupName,
			})
		}
	}

	return envelope, userId, groups, nil
}

// MainLoop starts the main program loop that will listen to all events on the previously initialized socket
// This method is blocking and will only exit if something goes horribly wrong
func MainLoop() {
	// Backend socket talks to workers over inproc
	backend, _ := zmq.NewSocket(zmq.DEALER)
	defer backend.Close()
	backend.Bind("inproc://backend")

	//  Launch pool of worker threads, precise number is not critical
	for i := 0; i < 5; i++ {
		go eventHandler()
	}

	//  Connect backend to frontend via a proxy
	err := zmq.Proxy(frontend, backend, nil)
	log.Fatalln("Proxy interrupted:", err)
}

func eventHandler() {
	worker, _ := zmq.NewSocket(zmq.DEALER)
	defer worker.Close()
	worker.Connect("inproc://backend")

	for {
		//  The DEALER socket gives us the reply envelope and message
		msg, _ := worker.RecvMessage(0)
		identity, content := pop(msg)



		envelope, userId, groups, err := getValidEvent(content[0])
		if err != nil {
			fmt.Println("failed to get event: ", err)
		}

		processEnvelope(envelope, userId, groups)

		payload := applyFilters(userId, envelope)
		message, _ := proto.Marshal(payload)
		worker.SendMessage(identity, message)
	}
}

func pop(msg []string) (head, tail []string) {
	if msg[1] == "" {
		head = msg[:2]
		tail = msg[2:]
	} else {
		head = msg[:1]
		tail = msg[1:]
	}
	return
}

func processEnvelope(envelope *Focus.Envelope, userId uuid.UUID, groups []models.Group) {
	db := database.Get()

	// Insert groups
	if len(groups) != 0 {
		if _, err := db.Model(&groups).OnConflict("DO NOTHING").Insert(); err != nil {
			fmt.Println("failed to insert group", err)
		}
	}

	for _, event := range envelope.Events {
		windowName := ""
		processName := ""
		afk := false

		if event.PayloadType == "Afk" {
			afk = true
		} else if event.PayloadType == "ContextChanged" {
			payload := &Focus.ContextEventPayload{}

			if err := ptypes.UnmarshalAny(event.Payload, payload); err != nil {
				fmt.Println("failed to unmarschal event", err)
				return
			}
			windowName = string(payload.WindowName)
			processName = string(payload.ProcessName)

		} else {
			return
		}


		deviceId, err := strconv.Atoi(envelope.DeviceID)
		if err != nil {
			fmt.Println(err)
			return // Couldn't read device ID, break
		}


		/////////////////////////////
		// INSERT EVERYTHING IN DB //
		/////////////////////////////


		// Insert event into DB
		eventToInsert := models.Event{
			UserId:      userId,
			DeviceId:    deviceId,
			WindowsName: windowName,
			ProcessName: processName,
			Afk:         afk,
			Time:        time.Unix(event.Timestamp.Seconds, int64(event.Timestamp.Nanos)),
		}
		if _, err := db.Model(&eventToInsert).Insert(); err != nil {
			fmt.Println("failed to insert event", err)
		}

		// Associate event with each group
		var groupEvents []models.GroupEvent
		for _, group := range groups {
			groupEvents = append(groupEvents, models.GroupEvent{
				EventId: eventToInsert.ID,
				GroupId: group.ID,
			})
		}
		if len(groups) != 0 {
			if _, err := db.Model(&groupEvents).Insert(); err != nil {
				fmt.Println("failed to associate event to group", err)
			}
		}

		fmt.Println("Inserted event for user", userId)
		fmt.Println("device: ", deviceId, "groups: ", groups)
		fmt.Println("window: ", windowName, "process: ", processName, "\n")
	}
}

func applyFilters(userId uuid.UUID, envelope *Focus.Envelope) *Focus.FilterEventPayload {
	enableDnd := false

	processName := ""
	afk := false
	event := envelope.Events[0]

	if event.PayloadType == "Afk" {
		afk = true
	} else if event.PayloadType == "ContextChanged" {
		payload := &Focus.ContextEventPayload{}

		if err := ptypes.UnmarshalAny(event.Payload, payload); err != nil {
			fmt.Println("failed to unmarschal event", err)
			return nil
		}
		processName = string(payload.ProcessName)

	} else {
		return nil
	}

	// Fetch user filters
	var filters []models.Filters
	err := database.Get().Model(&filters).Column("name").Where("user_id = ?", userId).Select()
	if err != nil {
		fmt.Println("failed to query filters", err)
	}

	for _, filter := range filters {
		if filter.Name == processName {
			fmt.Println("filter matched, activating DND")
			incrementDnd(userId)
			enableDnd = true
		} else if afk == true {
			enableDnd = false
		}
	}

	return &Focus.FilterEventPayload{IsDndOn:enableDnd}
}

func incrementDnd(userId uuid.UUID) {
	dnd := models.Dnd{
		UserId:      userId,
		Activations: 0,
	}

	_, err := database.Get().Model(&dnd).SelectOrInsert()
	if err != nil {
		fmt.Println("failed to update DND count", err)
	}

	dnd.Activations++
	_, err = database.Get().Model(&dnd).Where("user_id = ?", userId).Update()
	if err != nil {
		fmt.Println("failed to increment DND count", err)
	}
}