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


// getAllEvents will extract all query parameters for /windows and /processes routes.
//
func getAllEvents(request *routing.Context) ([]*models.Event, error) {
	// Extract parameters from query
	from := string(request.PostArgs().Peek("from"))
	to := string(request.PostArgs().Peek("to"))
	device := string(request.PostArgs().Peek("device"))
	group := string(request.PostArgs().Peek("group"))
	userId := request.Get("userId").(uuid.UUID)


	// Only one parameter allowed
	if device != "" && group != "" {
		request.Error("cannot use device and group parameters at the same time", fasthttp.StatusPreconditionFailed)
		return nil, errors.New("cannot use device and group parameters at the same time")
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
			return nil, err
		}
		fromTime := time.Unix(i, 0)
		query = query.Where("time > ?", fromTime)
	}
	if to != "" {
		i, err := strconv.ParseInt(from, 10, 64)
		if err != nil {
			request.Error(err.Error(), fasthttp.StatusInternalServerError)
			return nil, err
		}
		toTime := time.Unix(i, 0)
		query = query.Where("time < ?", toTime)
	}

	err := query.Select()
	if err != nil {
		fmt.Println("failed to query",  err)
		request.Error(err.Error(), fasthttp.StatusInternalServerError)
		return nil, err
	}

	return events, nil
}