package api

import (
	"encoding/json"
	"fmt"
	"github.com/FocusCompany/backend-go/database"
	"github.com/FocusCompany/backend-go/models"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

type Activity struct {
	Process string
	Window  string
	Start   time.Time
	End     time.Time
	Afk     bool
}

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
	// Extract parameters from query
	from := string(request.PostArgs().Peek("from"))
	to := string(request.PostArgs().Peek("to"))
	device := string(request.PostArgs().Peek("device"))
	group := string(request.PostArgs().Peek("group"))
	userId := request.Get("userId").(uuid.UUID)


	// Only one parameter allowed
	if device != "" && group != "" {
		request.Error("cannot use device and group parameters at the same time", fasthttp.StatusPreconditionFailed)
		return nil
	}

	// Build query
	var events []*models.Event
	query := database.Get().Model(&events).Where("user_id = ?", userId)

	if device != "" {
		query = query.Where("device_id = ?", device)
	}
	if group != "" {
		query = query.Where("group_id = ?", group)
	}
	if from != "" {
		i, err := strconv.ParseInt(from, 10, 64)
		if err != nil {
			request.Error(err.Error(), fasthttp.StatusInternalServerError)
			return nil
		}
		fromTime := time.Unix(i, 0)
		query = query.Where("time > ?", fromTime)
	}
	if to != "" {
		i, err := strconv.ParseInt(from, 10, 64)
		if err != nil {
			request.Error(err.Error(), fasthttp.StatusInternalServerError)
			return nil
		}
		toTime := time.Unix(i, 0)
		query = query.Where("time < ?", toTime)
	}

	err := query.Select()
	if err != nil {
		fmt.Println("failed to query",  err)
		request.Error(err.Error(), fasthttp.StatusInternalServerError)
		return nil
	}







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
			if windows[i].Window == currentActivity.Window && windows[i].End.Add(45 * time.Minute).After(currentActivity.Start) {
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
		request.Error("failed to marhsal JSON" + err.Error(), fasthttp.StatusInternalServerError)
	}
	request.SetBody(response)
	return nil
}
