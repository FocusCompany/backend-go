package api

import (
	"errors"
	"fmt"
	"github.com/FocusCompany/backend-go/database"
	"github.com/FocusCompany/backend-go/models"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"time"
)

// Activity represents an event in a Go struct.
type Activity struct {
	Process string `json:"process"`
	Window  string `json:"window"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Afk     bool `json:"afk"`
}

// queryParam is a struct that holds the common query parameters for the stat routes
type queryParam struct {
	From   time.Time
	To     time.Time
	Device int
	Group  int
	UserId uuid.UUID
}

// getQueryParam will extract all query parameters for /windows and /processes routes and put it in a struct.
func getQueryParam(request *routing.Context) (queryParam, error) {
	// Extract parameters from query
	from := string(request.PostArgs().Peek("from"))
	to := string(request.PostArgs().Peek("to"))
	deviceId := request.PostArgs().GetUintOrZero("device")
	groupId := request.PostArgs().GetUintOrZero("group")

	queryParam := queryParam{
		UserId: request.Get("userId").(uuid.UUID),
		Device: deviceId,
		Group: groupId,
	}

	// Only one parameter allowed
	if queryParam.Device != 0 && queryParam.Group != 0 {
		request.Error("cannot use device and group parameters at the same time", fasthttp.StatusPreconditionFailed)
		return queryParam, errors.New("cannot use device and group parameters at the same time")
	}

	if from != "" {
		i, err := time.Parse(time.RFC3339, from)
		if err != nil {
			request.Error(err.Error(), fasthttp.StatusInternalServerError)
			return queryParam, err
		}
		queryParam.From = i
	}
	if to != "" {
		i, err := time.Parse(time.RFC3339, to)
		if err != nil {
			request.Error(err.Error(), fasthttp.StatusInternalServerError)
			return queryParam, err
		}
		queryParam.To = i
	}

	return queryParam, nil
}

// getEvents will fetch all events based on parameters found in the request
func getEvents(request *routing.Context) ([]*models.Event, error) {
	param, err := getQueryParam(request)
	if err != nil {
		return nil, err
	}

	// Fetch all events from DB
	var events []*models.Event
	query := database.Get().Model(&events).Where("user_id = ?", param.UserId)

	fmt.Println(param.Group, param.Device, param.UserId, param.From, param.To)

	if param.Group != 0 { query = query.Where("group_id = ?", param.Group) }
	if param.Device != 0 { query = query.Where("device_id = ?", param.Device) }
	if !param.From.IsZero() { query = query.Where("time > ?", param.From) }
	if !param.To.IsZero() { query = query.Where("time < ?", param.To) }

	err = query.Select()
	if err != nil {
		fmt.Println("failed to query events",  err)
		request.Error(err.Error(), fasthttp.StatusInternalServerError)
		return nil, err
	}

	return events, nil
}