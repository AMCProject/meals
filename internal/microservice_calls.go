package internal

import (
	"bytes"
	"encoding/json"
	"meals/internal/config"
	"net/http"
)

type Endpoints struct {
}

type EndpointsI interface {
	GetUser(userId string) (user User, err error)
	GetCalendar(userId string, meal Meal, delete bool) (err error)
}

var httpClient = &http.Client{}

func (e *Endpoints) GetUser(userId string) (user User, err error) {
	request, err := http.NewRequest(http.MethodGet, config.Config.UsersURL+"user/"+userId, nil)
	if err != nil {
		return
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return User{}, err
	}
	if response.StatusCode > 299 {
		newError := new(ErrorResponse)
		err = json.NewDecoder(response.Body).Decode(&newError)
		return User{}, newError
	}
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return
}

func (e *Endpoints) GetCalendar(userId string, meal Meal, delete bool) (err error) {
	var calendar []Calendar
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
		newError := new(ErrorResponse)
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
		weekCalendar := UpdateWeekCalendar{From: c.Date, To: c.Date}
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
				newError := new(ErrorResponse)
				err = json.NewDecoder(response.Body).Decode(&newError)
				return newError
			}
		} else {
			data := Calendar{UserId: userId, MealId: meal.Id, Name: meal.Name, Date: c.Date}
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
				newError := new(ErrorResponse)
				err = json.NewDecoder(response.Body).Decode(&newError)
				return newError
			}
		}
	}

	return
}
