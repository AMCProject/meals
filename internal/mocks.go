package internal

import (
	"github.com/stretchr/testify/mock"
)

type EndpointsMock struct {
	mock.Mock
}

func (e *EndpointsMock) GetCalendar(userId string, meal Meal, delete bool) (err error) {
	args := e.Called(userId, meal, delete)
	return args.Error(0)
}

func (e *EndpointsMock) GetUser(userId string) (user User, err error) {
	args := e.Called(userId)
	return args.Get(0).(User), args.Error(1)
}
