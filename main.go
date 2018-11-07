package main

import (
	"fmt"
	"github.com/FocusCompany/backend-go/api"
	"github.com/FocusCompany/backend-go/database"
	"github.com/FocusCompany/backend-go/models"
	"github.com/pebbe/zmq4"
	"github.com/satori/go.uuid"
	"sync"
)

func main() {
	database.Init()

	go api.Init()

	zmq, err := zmq4.NewContext()
	if err != nil {
		fmt.Printf("failed to create ZMQ context")
	}

	_, err = zmq.NewSocket(zmq4.DEALER)
	if err != nil {
		fmt.Printf("failed to create socket")
	}

	userID := uuid.NewV4()
	groupID := uuid.NewV4()
	deviceID := uuid.NewV4()

	event := models.Event{
		UserId:      userID,
		GroupId:     groupID,
		DeviceId:    deviceID,
		WindowsName: "Skype",
		ProcessName: "/programfiles/Skype.exe",
	}

	db := database.Get()
	if _, err := db.Model(&event).Insert(); err != nil {
		fmt.Println("failed to insert event", err)
	}

	var getEvent models.Event
	db.Model(&getEvent).Where("id = ?", event.ID).Select()
	fmt.Printf(getEvent.ProcessName)

	// Block for other goroutines
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
