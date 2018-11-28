package main

import (
	"github.com/FocusCompany/backend-go/api"
	"github.com/FocusCompany/backend-go/database"
	"github.com/FocusCompany/backend-go/socket"
)

func main() {
	database.Init() // Create the connection to the DB
	go api.Init()   // Instanciate the API router

	err := socket.InitSocket() // Create a socket listening on tcp://*:5555
	if err != nil {
		panic(err)
	}

	socket.MainLoop() // Listen to message on socket
}
