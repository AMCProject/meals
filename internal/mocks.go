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
