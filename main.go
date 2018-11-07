package main

import (
	"fmt"
	"github.com/FocusCompany/backend-go/api"
	"github.com/FocusCompany/backend-go/database"
	"github.com/FocusCompany/backend-go/models"
	"github.com/satori/go.uuid"
)

func main() {
	database.Init()

	api.Init()

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
}
