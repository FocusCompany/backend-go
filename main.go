package main

import (
	"fmt"
	"github.com/FocusCompany/backend-go/api"
	"github.com/FocusCompany/backend-go/database"
	"github.com/FocusCompany/backend-go/models"
	"github.com/FocusCompany/backend-go/proto"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/pebbe/zmq4"
	"github.com/satori/go.uuid"
	"sync"
	"time"
)

func main() {
	database.Init()
	go api.Init()


	socket, err := zmq4.NewSocket(zmq4.PULL)
	defer socket.Close()
	if err != nil {
		fmt.Printf("failed to create socket")
		return
	}

	err = socket.Bind("tcp://*:5555")
	if err != nil {
		fmt.Printf("failed to bind socket")
		return
	}
	fmt.Println("Starting to receive")


	for ;; {
		received, err := socket.Recv(0)
		if err != nil {
			fmt.Printf("socket.recv error: %v", err)
			return
		}
		envelope := &Focus.Envelope{}
		err = proto.Unmarshal([]byte(received), envelope)
		if err != nil {
			fmt.Printf("got error: %v", err)
			return
		}

		// Check JWT from envelope
		userId, err := api.ValidateJwt(envelope.Jwt)
		if err != nil {
			fmt.Println("invalid JWT", err)
			return
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

			// Insert event into DB
			eventToInsert := models.Event{
				UserId:      userId,
				GroupId:     uuid.UUID{},
				DeviceId:    uuid.FromStringOrNil(envelope.DeviceID),
				WindowsName: windowName,
				ProcessName: processName,
				Afk:         afk,
				Time:        time.Unix(event.Timestamp.Seconds, int64(event.Timestamp.Nanos)),
			}

			db := database.Get()
			if _, err := db.Model(&eventToInsert).Insert(); err != nil {
				fmt.Println("failed to insert event", err)
			}
			fmt.Println("Insert event", eventToInsert)
		}
	}

	// Block for other goroutines
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
