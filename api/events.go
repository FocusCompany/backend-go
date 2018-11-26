package api

import (
	"errors"
	"fmt"
	"github.com/FocusCompany/backend-go/database"
	"github.com/FocusCompany/backend-go/models"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"strconv"
	"time"
)

// Activity represents an event in a Go struct.
type Activity struct {
	Process string
	Window  string
	Start   time.Time
	End     time.Time
	Afk     bool
}

// queryParam is a struct that holds the common query parameters for the stat routes
type queryParam struct {
	From   time.Time
	To     time.Time
	Device string
	Group  string
	UserId uuid.UUID
}

// getQueryParam will extract all query parameters for /windows and /processes routes and put it in a struct.
func getQueryParam(request *routing.Context) (queryParam, error) {
	// Extract parameters from query
	from := string(request.PostArgs().Peek("from"))
	to := string(request.PostArgs().Peek("to"))

	queryParam := queryParam{
		Device: string(request.PostArgs().Peek("device")),
		Group:  string(request.PostArgs().Peek("group")),
		UserId: request.Get("userId").(uuid.UUID),
	}

	// Only one parameter allowed
	if queryParam.Device != "" && queryParam.Group != "" {
		request.Error("cannot use device and group parameters at the same time", fasthttp.StatusPreconditionFailed)
		return queryParam, errors.New("cannot use device and group parameters at the same time")
	}

	if from != "" {
		i, err := strconv.ParseInt(from, 10, 64)
		if err != nil {
			request.Error(err.Error(), fasthttp.StatusInternalServerError)
			return queryParam, err
		}
		queryParam.From = time.Unix(i, 0)
	}
	if to != "" {
		i, err := strconv.ParseInt(from, 10, 64)
		if err != nil {
			request.Error(err.Error(), fasthttp.StatusInternalServerError)
			return queryParam, err
		}
		queryParam.To = time.Unix(i, 0)
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

	if param.Group != "" { query = query.Where("group_id = ?", param.Group) }
	if param.Device != "" { query = query.Where("device_id = ?", param.Device) }
	if !param.From.IsZero() { query = query.Where("time > ?", param.From) }
	if !param.To.IsZero() { query = query.Where("time < ?", param.To) }

	err = query.Select()
	if err != nil {
		fmt.Println("failed to query",  err)
		request.Error(err.Error(), fasthttp.StatusInternalServerError)
		return nil, err
	}

	return events, nil
}