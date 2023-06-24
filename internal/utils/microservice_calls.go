package utils

import (
	"bytes"
	"encoding/json"
	"meals/internal"
	"meals/internal/config"
	"meals/internal/models"
	"net/http"
)

type Endpoints struct {
}

type EndpointsI interface {
	GetCalendar(userId string, meal models.Meal, delete bool) (err error)
}

var httpClient = &http.Client{}

func (e *Endpoints) GetCalendar(userId string, meal models.Meal, delete bool) (err error) {
	var calendar []models.Calendar
	request, err := http.NewRequest(http.MethodGet, config.Config.CalendarsURL+"user/"+userId+"/calendar", nil)
	if err != nil {
		return
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode == 404 {
		return
	}
	if response.StatusCode > 299 {
		newError := new(internal.ErrorResponse)
		err = json.NewDecoder(response.Body).Decode(&newError)
		return newError
	}
	err = json.NewDecoder(response.Body).Decode(&calendar)
	if err != nil {
		return err
	}
	for _, c := range calendar {
		if c.MealId != meal.Id {
			continue
		}
		weekCalendar := models.UpdateWeekCalendar{From: c.Date, To: c.Date}
		if delete {
			content, _ := json.Marshal(weekCalendar)
			request, err = http.NewRequest(http.MethodPut, config.Config.CalendarsURL+"user/"+userId+"/redoweek", bytes.NewBuffer(content))
			request.Header.Set("Content-Type", "application/json;charset=UTF-8")
			if err != nil {
				return
			}
			response, err = httpClient.Do(request)
			if err != nil {
				return err
			}
			if response.StatusCode > 299 {
				newError := new(internal.ErrorResponse)
				err = json.NewDecoder(response.Body).Decode(&newError)
				return newError
			}
		} else {
			data := models.Calendar{UserId: userId, MealId: meal.Id, Name: meal.Name, Date: c.Date}
			content, _ := json.Marshal(data)
			request, err = http.NewRequest(http.MethodPut, config.Config.CalendarsURL+"user/"+userId+"/calendar", bytes.NewBuffer(content))
			if err != nil {
				return
			}
			request.Header.Set("Content-Type", "application/json;charset=UTF-8")
			response, err = httpClient.Do(request)
			if err != nil {
				return err
			}
			if response.StatusCode > 299 {
				newError := new(internal.ErrorResponse)
				err = json.NewDecoder(response.Body).Decode(&newError)
				return newError
			}
		}
	}

	return
}
