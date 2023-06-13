package internal

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	RouteMeal       = "/user/:user_id/meal"
	RouteMealID     = "/user/:user_id/meal/:id"
	RouteMealOthers = "/user/:user_id/others"

	RouteIngredients = "/ingredients"

	ParamUserID = "user_id"
	ParamMealID = "id"
)

type ErrorResponse struct {
	Err ErrorBody `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Err.Status, e.Err.Message)
}

type ErrorBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewErrorResponse(c echo.Context, err error) error {
	errResponse := &ErrorResponse{Err: errorsMap[err.Error()]}
	if errResponse.Err.Status == 0 {
		if err := c.JSON(http.StatusInternalServerError, err); err != nil {
			return err
		}
		return err
	}
	if err := c.JSON(errResponse.Err.Status, errResponse); err != nil {
		return err
	}
	return errResponse
}

var errorsMap = map[string]ErrorBody{
	ErrUserIDNotPresent.Error():   {Status: http.StatusBadRequest, Message: ErrUserIDNotPresent.Error()},
	ErrMealIDNotPresent.Error():   {Status: http.StatusBadRequest, Message: ErrMealIDNotPresent.Error()},
	ErrMealTypeNotPresent.Error(): {Status: http.StatusBadRequest, Message: ErrMealTypeNotPresent.Error()},
	ErrSomethingWentWrong.Error(): {Status: http.StatusInternalServerError, Message: ErrSomethingWentWrong.Error()},
	ErrWrongBody.Error():          {Status: http.StatusBadRequest, Message: ErrWrongBody.Error()},
	ErrMealNotFound.Error():       {Status: http.StatusNotFound, Message: ErrMealNotFound.Error()},
	ErrMealsNotFound.Error():      {Status: http.StatusNotFound, Message: ErrMealsNotFound.Error()},
	ErrUserNotFound.Error():       {Status: http.StatusNotFound, Message: ErrUserNotFound.Error()},
}
var (
	ErrUserIDNotPresent   = errors.New("error with userID given")
	ErrMealIDNotPresent   = errors.New("error with mealID given")
	ErrMealTypeNotPresent = errors.New("error with mealType given")
	ErrSomethingWentWrong = errors.New("something went wrong")
	ErrWrongBody          = errors.New("malformed body")
	ErrMealNotFound       = errors.New("meal not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrMealsNotFound      = errors.New("meals not found")
)
