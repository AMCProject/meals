package internal

import (
	"encoding/json"
	"meals/internal/config"
	"net/http"
)

type Endpoints struct {
}

type EndpointsI interface {
	GetUser(userId string) (user User, err error)
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
	err = json.NewDecoder(response.Body).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return
}
