package api

import (
	"encoding/json"
	"fmt"
	"github.com/FocusCompany/backend-go/database"
	"github.com/FocusCompany/backend-go/models"
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
func windowHandler(request *routing.Context) error {
	events, err := getEvents(request)

	// Compute time per Process
	var windows []Activity
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
		for i := len(windows) - 1; i >= 0; i-- {

			//If the two events are more than 45 minutes apart
			if windows[i].Window == currentActivity.Window && windows[i].End.Add(10 * time.Minute).After(currentActivity.Start) {
				windows[i].End = currentActivity.End
				isNew = false
			}
		}

		if isNew && !currentActivity.Start.IsZero() {
			windows = append(windows, currentActivity)
		}
		currentActivity = newActivity
	}

	response, err := json.Marshal(windows)
	if err != nil {
		request.Error("failed to marshal JSON" + err.Error(), fasthttp.StatusInternalServerError)
	}
	request.SetBody(response)
	return nil
}

func windowListHandler(request *routing.Context) error {
	param, err := getQueryParam(request)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Fetch all events from DB
	var events []*models.Event
	query := database.Get().
		Model(&events).
		Where("user_id = ?", param.UserId).
		ColumnExpr("DISTINCT window_name")

	if param.Group != 0 { query = query.Where("group_id = ?", param.Group) }
	if param.Device != 0 { query = query.Where("device_id = ?", param.Device) }
	if !param.From.IsZero() { query = query.Where("time > ?", param.From) }
	if !param.To.IsZero() { query = query.Where("time < ?", param.To) }

	count, err := query.SelectAndCount()
	if err != nil {
		fmt.Println("failed to query",  err)
		request.Error(err.Error(), fasthttp.StatusInternalServerError)
		return err
	}

	windows := make([]string, count)
	for i, event := range events {
		windows[i] = event.WindowsName
	}

	response, err := json.Marshal(windows)
	if err != nil {
		request.Error("failed to marshal JSON" + err.Error(), fasthttp.StatusInternalServerError)
	}
	request.SetBody(response)
	return nil
}
