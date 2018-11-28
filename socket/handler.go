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

func getValidEvent(received string) (*Focus.Envelope, uuid.UUID, error) {
	envelope := &Focus.Envelope{}
	err := proto.Unmarshal([]byte(received), envelope)
	if err != nil {
		return nil, uuid.UUID{}, errors.New("proto.Unmarshal error " + err.Error())
	}

	// Check JWT from envelope
	userId, err := api.ValidateJwt(envelope.Jwt)
	if err != nil {
		return nil, uuid.UUID{}, errors.New("Invalid JWT " + err.Error())
	}
	return envelope, userId, nil
}

// MainLoop starts the main program loop that will listen to all events on the previously initialized socket
// This method is blocking and will only exit if something goes horribly wrong
func MainLoop(sock *zmq.Socket) {
	//  Backend socket talks to workers over inproc
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

		worker.SendMessage(identity, "true")

		envelope, userId, err := getValidEvent(content[0])
		if err != nil {
			fmt.Println("failed to get event: ", err)
		}

		processEnvelope(envelope, userId)
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

func processEnvelope(envelope *Focus.Envelope, userId uuid.UUID) {
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

		// Insert event into DB
		eventToInsert := models.Event{
			UserId:      userId,
			GroupId:     0,
			DeviceId:    deviceId,
			WindowsName: windowName,
			ProcessName: processName,
			Afk:         afk,
			Time:        time.Unix(event.Timestamp.Seconds, int64(event.Timestamp.Nanos)),
		}

		db := database.Get()
		if _, err := db.Model(&eventToInsert).Insert(); err != nil {
			fmt.Println("failed to insert event", err)
		}

		fmt.Println("Inserted event for user", userId)
		fmt.Println("device: ", deviceId, "group: ", 0)
		fmt.Println("window: ", windowName, "process: ", processName, "\n")
	}
}
