package api

import (
	"encoding/json"
	"github.com/FocusCompany/backend-go/database"
	"github.com/FocusCompany/backend-go/models"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
)

func getTotalDndHandler(request *routing.Context) error {
	var dnd models.Dnd
	database.Get().Model(&dnd).Where("user_id = ?", request.Get("userId").(uuid.UUID)).Select()

	response, err := json.Marshal(dnd)
	if err != nil {
		request.Error("failed to marshal JSON" + err.Error(), fasthttp.StatusInternalServerError)
	}
	request.SetBody(response)
	return nil
}
