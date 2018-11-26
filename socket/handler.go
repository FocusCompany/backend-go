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
	"github.com/pebbe/zmq4"
	"github.com/satori/go.uuid"
	"time"
)

func getValidEvent(sock *zmq4.Socket) (*Focus.Envelope, uuid.UUID, error) {
	received, err := sock.Recv(0)
	if err != nil {
		return nil, uuid.UUID{}, errors.New("socket.recv error " + err.Error())
	}
	envelope := &Focus.Envelope{}
	err = proto.Unmarshal([]byte(received), envelope)
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
func MainLoop(sock *zmq4.Socket) {
	for ;; {
		envelope, userId, err := getValidEvent(sock)
		if err != nil {
			fmt.Println("FAILED TO RECEIVE EVENT", err)
		}


		// Insert events in DB
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
			fmt.Println("Inserted event", eventToInsert)
		}
	}

}
