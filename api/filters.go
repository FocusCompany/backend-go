package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/FocusCompany/backend-go/database"
	"github.com/FocusCompany/backend-go/models"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"strings"
)

func getFiltersHandler(request *routing.Context) error {
	userId := request.Get("userId").(uuid.UUID)

	filterList, err := GetFilters(userId)
	if err != nil {
		request.Error(err.Error(), fasthttp.StatusInternalServerError)
		return err
	}

	response, err := json.Marshal(filterList)
	if err != nil {
		request.Error("failed to marshal JSON" + err.Error(), fasthttp.StatusInternalServerError)
	}
	request.SetBody(response)
	return nil
}

// GetFilters will return all filters associated with a given userId
func GetFilters(userId uuid.UUID) ([]string, error) {
	var filters []models.Filters
	count, err := database.Get().
		Model(&filters).
		ColumnExpr("DISTINCT name").
		Where("user_id = ?", userId).
		SelectAndCount()

	if err != nil {
		fmt.Println("getFiltersHandler " + err.Error())
		return nil, err
	}

	filterList := make([]string, count)
	for i, filter := range filters {
		filterList[i] = filter.Name
	}

	return filterList, nil
}

func updateFiltersHandler(request *routing.Context) error {
	userId := request.Get("userId").(uuid.UUID)
	newFilters := strings.Split(string(request.PostArgs().Peek("filters")), ",")

	if len(newFilters) == 0 {
		request.Error("missing filter list", fasthttp.StatusFailedDependency)
		return errors.New("missing filter list")
	}

	// Get new filters from request
	var filters []models.Filters
	for _, filter := range newFilters {
		filters = append(filters, models.Filters{
			UserId: userId,
			Name:   filter,
		})
	}

	// Delete all previous filters
	query, err := database.Get().Begin() // Prepare SQL transaction
	_, err = query.Model((*models.Filters)(nil)).
		Where("user_id = ?", userId).
		Delete()
	if err != nil {
		query.Rollback()
		fmt.Println(err.Error())
		request.Error("failed to delete previous filters", fasthttp.StatusInternalServerError)
		return err
	}

	// Insert new ones
	_, err = query.Model(&filters).Insert()
	if err != nil {
		query.Rollback()
		fmt.Println(err.Error())
		request.Error("failed to insert new filters", fasthttp.StatusInternalServerError)
		return err
	}

	// Commit transaction
	err = query.Commit()
	if err != nil {
		fmt.Println(err.Error())
		request.Error("failed to update filters", fasthttp.StatusInternalServerError)
		return err
	}

	return nil
}