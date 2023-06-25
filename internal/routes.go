package internal

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	RouteMeal          = "/user/:user_id/meal"
	RouteMealID        = "/user/:user_id/meal/:id"
	RouteExternalMeals = "/meals"

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
	ErrorWithExternalAPI.Error():  {Status: http.StatusInternalServerError, Message: ErrorWithExternalAPI.Error()},
	ErrWrongBody.Error():          {Status: http.StatusBadRequest, Message: ErrWrongBody.Error()},
	ErrMealNotFound.Error():       {Status: http.StatusNotFound, Message: ErrMealNotFound.Error()},
	ErrMealsNotFound.Error():      {Status: http.StatusNotFound, Message: ErrMealsNotFound.Error()},
	ErrUserNotFound.Error():       {Status: http.StatusNotFound, Message: ErrUserNotFound.Error()},
	ErrMealAlreadyExist.Error():   {Status: http.StatusConflict, Message: ErrMealAlreadyExist.Error()},
}
var (
	ErrUserIDNotPresent   = errors.New("error con el ID de usuario indicado")
	ErrMealIDNotPresent   = errors.New("error con el ID de comida indicado")
	ErrMealTypeNotPresent = errors.New("error con el tipo de comida indicado")
	ErrSomethingWentWrong = errors.New("error inesperado")
	ErrWrongBody          = errors.New("el cuerpo enviado es erróneo")
	ErrMealNotFound       = errors.New("comida no encontrada")
	ErrUserNotFound       = errors.New("usuario no encontrado")
	ErrMealsNotFound      = errors.New("comidas no encontradas")
	ErrorWithExternalAPI  = errors.New("error inesperado con la API externa")
	ErrMealAlreadyExist   = errors.New("ya existe una comida con este nombre")
)
