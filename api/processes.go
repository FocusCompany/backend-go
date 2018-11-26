package api

import (
	"encoding/json"
	"fmt"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"time"
)

// PARAMS
//
// one of
//	- device:	id		- ID of device
//	- group:	id		- ID of groupe
//  - empty:			- All user devices
//
// optional
// 	- from:		date	- Beginning of events to get
// 	- to:		date	- End of events to get
func processHandler(request *routing.Context) error {
	events, err := getAllEvents(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Compute time per Process
	var processes []Activity
	var currentActivity Activity
	for _, event := range events {
		newActivity := Activity{
			Process: event.ProcessName,
			Window:  event.WindowsName,
			Start:   event.Time,
			Afk:     event.Afk,
		}

		currentActivity.End = newActivity.Start

		// Get the last activity with the same name as the current activity
		isNew := true
		for i := len(processes) - 1; i >= 0; i-- {

			//If the two events are more than 45 minutes apart
			if processes[i].Process == currentActivity.Process && processes[i].End.Add(10 * time.Minute).After(currentActivity.Start) {
				processes[i].End = currentActivity.End
				isNew = false
			}
		}

		if isNew && !currentActivity.Start.IsZero() {
			processes = append(processes, currentActivity)
		}
		currentActivity = newActivity
	}




	response, err := json.Marshal(processes)
	if err != nil {
		request.Error("failed to marshal JSON" + err.Error(), fasthttp.StatusInternalServerError)
	}
	request.SetBody(response)
	return nil
}

