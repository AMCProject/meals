package internal

import (
	"encoding/json"
	"net/http"
)

const userUrl = "http://localhost:3100/"

type Endpoints struct {
}

type EndpointsI interface {
	GetUser(userId string) (user User, err error)
}

var httpClient = &http.Client{}

func (e *Endpoints) GetUser(userId string) (user User, err error) {
	request, err := http.NewRequest(http.MethodGet, userUrl+"user/"+userId, nil)
	if err != nil {
		return
	}
	response, err := httpClient.Do(request)
	if err != nil {
		return User{}, err
	}
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return
}
