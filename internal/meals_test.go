package internal

import (
	"bytes"
	"github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/suite"
	"meals/pkg/database"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var databaseTest = "/amc_test.db"

type MealAPITestSuite struct {
	suite.Suite
	db       *database.Database
	httpMock *EndpointsMock
}

func TestMealAPITestSuite(t *testing.T) {
	suite.Run(t, new(MealAPITestSuite))
}

func (s *MealAPITestSuite) SetupTest() {
	s.httpMock = &EndpointsMock{}
	Microservices = s.httpMock
	_ = database.RemoveDB(databaseTest)
	s.db = database.InitDB(databaseTest)
	s.db.Conn.Exec(createMeal, "01FN3EEB2NVFJAHAPM00000001", "01FN3EEB2NVFJAHAPU00000001", "pizza", "", "", "ocasional", "Tomate,Queso,Pollo", 130, "invierno,verano")
	s.db.Conn.Exec(createMeal, "01FN3EEB2NVFJAHAPM00000002", "01FN3EEB2NVFJAHAPU00000001", "ensalada", "", "", "semanal", "Tomate,Lechuga,Cebolla,Aguacate", 100, "general")
}

func (s *MealAPITestSuite) TearDownTest() {
	s.db = nil
	_ = database.RemoveDB(databaseTest)
}

func (s *MealAPITestSuite) TestPostMealHandler() {
	tests := []struct {
		name               string
		reqBody            interface{}
		userId             string
		expectedULID       ulid.ULID
		expectedResp       interface{}
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:   "[001] Create new meal (ok)",
			userId: "01FN3EEB2NVFJAHAPU00000001",
			reqBody: &Meal{
				Name:        "Huevos fritos con patatas",
				Description: "",
				Image:       "",
				Type:        "ocasional",
				Ingredients: []string{"Patata frita", "Huevo frito"},
				Seasons:     []string{"general"},
			},
			expectedULID: ulid.MustParse("01FN3EEB2NVFJAHAPM00000003"),
			expectedResp: &Meal{
				Id:          "01FN3EEB2NVFJAHAPM00000003",
				Name:        "Huevos fritos con patatas",
				Description: "",
				Image:       "",
				Type:        "ocasional",
				Ingredients: []string{"Patata frita", "Huevo frito"},
				Kcal:        208,
				Seasons:     []string{"general"},
			},
			expectedStatusCode: http.StatusCreated,
			wantErr:            false,
		},
		{
			name:   "[002] Wrong meal struct, name is missing (400)",
			userId: "01FN3EEB2NVFJAHAPU00000001",
			reqBody: &Meal{
				UserId:      "01FN3EEB2NVFJAHAPU00000004",
				Description: "",
				Image:       "",
				Type:        "ocasional",
				Ingredients: []string{"Patata frita, Huevo frito"},
				Kcal:        320,
				Seasons:     []string{"general"},
			},
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusBadRequest,
					Message: ErrWrongBody.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
		{
			name:    "[003] Wrong struct sent (400)",
			userId:  "01FN3EEB2NVFJAHAPU00000001",
			reqBody: "invalid",
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusBadRequest,
					Message: ErrWrongBody.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
		{
			name: "[004] User id not present (400)",
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusBadRequest,
					Message: ErrUserIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
	}
	getEchoContext := func(userId string, request interface{}) echo.Context {
		var body []byte
		body, err := jsoniter.Marshal(request)
		s.NoError(err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, RouteMeal, bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(ParamUserID)
		c.SetParamValues(userId)
		return c
	}

	for _, t := range tests {
		s.Run(t.name, func() {

			userManager := NewMealManager(*s.db)
			api := MealAPI{DB: *s.db, Manager: userManager}

			s.httpMock.On("GetUser", t.userId).Return(User{}, nil).Once()

			c := getEchoContext(t.userId, t.reqBody)
			err := api.PostMealHandler(c)

			if t.wantErr {
				s.Equal(t.wantErr, err != nil)
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				errorReturned := new(ErrorResponse)
				s.NoError(jsoniter.Unmarshal(body, errorReturned))
				s.Equal(errorReturned, t.expectedResp)
			} else {
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				actualMeal := new(Meal)
				s.NoError(jsoniter.Unmarshal(body, actualMeal))
				actualMeal.Id = t.expectedULID.String()
				s.Equal(actualMeal, t.expectedResp)
			}

			s.Equal(t.expectedStatusCode, c.Response().Status)
		})
	}
}

func (s *MealAPITestSuite) TestGetMealHandler() {
	tests := []struct {
		name               string
		userID             string
		mealID             string
		expectedResp       interface{}
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:   "[001] Get meal (ok)",
			userID: "01FN3EEB2NVFJAHAPU00000001",
			mealID: "01FN3EEB2NVFJAHAPM00000001",
			expectedResp: &Meal{
				Id:          "01FN3EEB2NVFJAHAPM00000001",
				Name:        "pizza",
				Description: "",
				Image:       "",
				Type:        "ocasional",
				Ingredients: []string{"Tomate", "Queso", "Pollo"},
				Kcal:        130,
				Seasons:     []string{"invierno", "verano"},
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name:   "[002] Get meal, userId not indicated (400)",
			mealID: "01FN3EEB2NVFJAHAPM00000001",
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusBadRequest,
					Message: ErrUserIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
		{
			name:   "[003] Get meal, mealId not indicated (400)",
			userID: "01FN3EEB2NVFJAHAPU00000001",
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusBadRequest,
					Message: ErrMealIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
		{
			name:   "[005] Meal does not exist (404)",
			userID: "01FN3EEB2NVFJAHAPU00000001",
			mealID: "01FN3EEB2NVFJAHAPM00000099",
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusNotFound,
					Message: ErrMealNotFound.Error(),
				},
			},
			expectedStatusCode: http.StatusNotFound,
			wantErr:            true,
		},
	}
	getEchoContext := func(userId, mealId string) echo.Context {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, RouteMealID, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(ParamUserID, ParamMealID)
		c.SetParamValues(userId, mealId)
		return c
	}
	for _, t := range tests {
		s.Run(t.name, func() {
			userManager := NewMealManager(*s.db)
			api := MealAPI{DB: *s.db, Manager: userManager}
			s.httpMock.On("GetUser", t.userID).Return(User{}, nil).Once()

			c := getEchoContext(t.userID, t.mealID)
			err := api.GetMealHandler(c)

			if t.wantErr {
				s.Equal(t.wantErr, err != nil)
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				errorReturned := new(ErrorResponse)
				s.NoError(jsoniter.Unmarshal(body, errorReturned))
				s.Equal(errorReturned, t.expectedResp)
			} else {
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				actualMeal := new(Meal)
				s.NoError(jsoniter.Unmarshal(body, actualMeal))
				s.Equal(actualMeal, t.expectedResp)
			}

			s.Equal(t.expectedStatusCode, c.Response().Status)
		})
	}
}

func (s *MealAPITestSuite) TestListMealsHandler() {
	tests := []struct {
		name               string
		userID             string
		mealID             string
		filters            url.Values
		expectedResp       interface{}
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:   "List meals (ok)",
			userID: "01FN3EEB2NVFJAHAPU00000001",
			expectedResp: &[]Meal{
				{
					Id:          "01FN3EEB2NVFJAHAPM00000001",
					Name:        "pizza",
					Description: "",
					Image:       "",
					Type:        "ocasional",
					Ingredients: []string{"Tomate", "Queso", "Pollo"},
					Kcal:        130,
					Seasons:     []string{"invierno", "verano"},
				},
				{
					Id:          "01FN3EEB2NVFJAHAPM00000002",
					Name:        "ensalada",
					Description: "",
					Image:       "",
					Type:        "semanal",
					Ingredients: []string{"Tomate", "Lechuga", "Cebolla", "Aguacate"},
					Kcal:        100,
					Seasons:     []string{"general"},
				},
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "List meals filtered by name (ok)",
			filters: map[string][]string{
				"name": {"pizza"},
			},
			userID: "01FN3EEB2NVFJAHAPU00000001",
			expectedResp: &[]Meal{
				{
					Id:          "01FN3EEB2NVFJAHAPM00000001",
					Name:        "pizza",
					Description: "",
					Image:       "",
					Type:        "ocasional",
					Ingredients: []string{"Tomate", "Queso", "Pollo"},
					Kcal:        130,
					Seasons:     []string{"invierno", "verano"},
				},
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "List meals filtered by type (ok)",
			filters: map[string][]string{
				"type": {"semanal"},
			},
			userID: "01FN3EEB2NVFJAHAPU00000001",
			expectedResp: &[]Meal{
				{
					Id:          "01FN3EEB2NVFJAHAPM00000002",
					Name:        "ensalada",
					Description: "",
					Image:       "",
					Type:        "semanal",
					Ingredients: []string{"Tomate", "Lechuga", "Cebolla", "Aguacate"},
					Kcal:        100,
					Seasons:     []string{"general"},
				},
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "List meals filtered by healthy (ok)",
			filters: map[string][]string{
				"healthy": {"true"},
			},
			userID: "01FN3EEB2NVFJAHAPU00000001",
			expectedResp: &[]Meal{
				{
					Id:          "01FN3EEB2NVFJAHAPM00000002",
					Name:        "ensalada",
					Description: "",
					Image:       "",
					Type:        "semanal",
					Ingredients: []string{"Tomate", "Lechuga", "Cebolla", "Aguacate"},
					Kcal:        100,
					Seasons:     []string{"general"},
				},
				{
					Id:          "01FN3EEB2NVFJAHAPM00000001",
					Name:        "pizza",
					Description: "",
					Image:       "",
					Type:        "ocasional",
					Ingredients: []string{"Tomate", "Queso", "Pollo"},
					Kcal:        130,
					Seasons:     []string{"invierno", "verano"},
				},
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "List meals filtered by seasons (ok)",
			filters: map[string][]string{
				"[]season": {"invierno"},
			},
			userID: "01FN3EEB2NVFJAHAPU00000001",
			expectedResp: &[]Meal{
				{
					Id:          "01FN3EEB2NVFJAHAPM00000001",
					Name:        "pizza",
					Description: "",
					Image:       "",
					Type:        "ocasional",
					Ingredients: []string{"Tomate", "Queso", "Pollo"},
					Kcal:        130,
					Seasons:     []string{"invierno", "verano"},
				},
				{
					Id:          "01FN3EEB2NVFJAHAPM00000002",
					Name:        "ensalada",
					Description: "",
					Image:       "",
					Type:        "semanal",
					Ingredients: []string{"Tomate", "Lechuga", "Cebolla", "Aguacate"},
					Kcal:        100,
					Seasons:     []string{"general"},
				},
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name: "List meals, userId not indicated (400)",
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusBadRequest,
					Message: ErrUserIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
	}
	getEchoContext := func(userId string, filters url.Values) echo.Context {
		e := echo.New()
		queryString := ""
		if len(filters) > 0 {
			queryString = "/?" + filters.Encode()
		}
		req := httptest.NewRequest(http.MethodGet, RouteMeal+queryString, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(ParamUserID)
		c.SetParamValues(userId)

		return c
	}
	for _, t := range tests {
		s.Run(t.name, func() {
			userManager := NewMealManager(*s.db)
			api := MealAPI{DB: *s.db, Manager: userManager}

			s.httpMock.On("GetUser", t.userID).Return(User{}, nil).Once()

			c := getEchoContext(t.userID, t.filters)
			err := api.ListMealsHandler(c)

			if t.wantErr {
				s.Equal(t.wantErr, err != nil)
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				errorReturned := new(ErrorResponse)
				s.NoError(jsoniter.Unmarshal(body, errorReturned))
				s.Equal(errorReturned, t.expectedResp)
			} else {
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				actualMeals := new([]Meal)
				s.NoError(jsoniter.Unmarshal(body, actualMeals))
				s.Equal(actualMeals, t.expectedResp)
			}

			s.Equal(t.expectedStatusCode, c.Response().Status)
		})
	}
}

func (s *MealAPITestSuite) TestPutMealHandler() {
	tests := []struct {
		name               string
		userID             string
		mealID             string
		reqBody            interface{}
		expectedResp       interface{}
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:   "Update meal (ok)",
			userID: "01FN3EEB2NVFJAHAPU00000001",
			mealID: "01FN3EEB2NVFJAHAPM00000001",
			reqBody: &Meal{
				Name:        "pizza margarita",
				Description: "",
				Image:       "",
				Type:        "ocasional",
				Ingredients: []string{"Tomate", "Queso"},
				Seasons:     []string{"invierno", "verano"},
			},
			expectedResp: &Meal{
				Id:          "01FN3EEB2NVFJAHAPM00000001",
				Name:        "pizza margarita",
				Description: "",
				Image:       "",
				Type:        "ocasional",
				Ingredients: []string{"Tomate", "Queso"},
				Kcal:        0,
				Seasons:     []string{"invierno", "verano"},
			},
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},
		{
			name:   "Update meal, userId not indicated (400)",
			mealID: "01FN3EEB2NVFJAHAPM00000001",
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusBadRequest,
					Message: ErrUserIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
		{
			name:   "Update meal, mealId not indicated (400)",
			userID: "01FN3EEB2NVFJAHAPU00000001",
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusBadRequest,
					Message: ErrMealIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
		{
			name:   "Update meal that does not exist (404)",
			userID: "01FN3EEB2NVFJAHAPU00000001",
			mealID: "01FN3EEB2NVFJAHAPM00000099",
			reqBody: &Meal{
				Id:          "01FN3EEB2NVFJAHAPM00000099",
				UserId:      "01FN3EEB2NVFJAHAPU00000001",
				Name:        "invented food",
				Description: "",
				Image:       "",
				Type:        "ocasional",
				Ingredients: []string{"Tomate", "Queso"},
				Kcal:        100,
				Seasons:     []string{"invierno"},
			},
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusNotFound,
					Message: ErrMealNotFound.Error(),
				},
			},
			expectedStatusCode: http.StatusNotFound,
			wantErr:            true,
		},
		{
			name:    "Update meal wrong body (400)",
			userID:  "01FN3EEB2NVFJAHAPU00000001",
			mealID:  "01FN3EEB2NVFJAHAPM00000001",
			reqBody: "invalid",
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusBadRequest,
					Message: ErrWrongBody.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
	}
	getEchoContext := func(userId, mealId string, request interface{}) echo.Context {
		var body []byte
		body, err := jsoniter.Marshal(request)
		s.NoError(err)
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, RouteMealID, bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(ParamUserID, ParamMealID)
		c.SetParamValues(userId, mealId)
		return c
	}
	for _, t := range tests {
		s.Run(t.name, func() {
			userManager := NewMealManager(*s.db)
			api := MealAPI{DB: *s.db, Manager: userManager}
			if _, ok := t.reqBody.(*Meal); ok {
				meal := t.reqBody.(*Meal)
				meal.Id = t.mealID
				s.httpMock.On("GetCalendar", t.userID, *meal, false).Return(nil).Once()
			}
			s.httpMock.On("GetUser", t.userID).Return(User{}, nil).Once()

			c := getEchoContext(t.userID, t.mealID, t.reqBody)
			err := api.PutMealHandler(c)

			if t.wantErr {
				s.Equal(t.wantErr, err != nil)
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				errorReturned := new(ErrorResponse)
				s.NoError(jsoniter.Unmarshal(body, errorReturned))
				s.Equal(errorReturned, t.expectedResp)
			} else {
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				actualMeal := new(Meal)
				s.NoError(jsoniter.Unmarshal(body, actualMeal))
				s.Equal(actualMeal, t.expectedResp)
			}

			s.Equal(t.expectedStatusCode, c.Response().Status)
		})
	}
}

func (s *MealAPITestSuite) TestDeleteMealHandler() {
	tests := []struct {
		name               string
		userID             string
		mealID             string
		expectedResp       interface{}
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:               "[001] Delete meal (ok)",
			userID:             "01FN3EEB2NVFJAHAPU00000001",
			mealID:             "01FN3EEB2NVFJAHAPM00000001",
			expectedStatusCode: http.StatusNoContent,
			wantErr:            false,
		},
		{
			name:   "[002] Delete meal, userId not indicated (400)",
			mealID: "01FN3EEB2NVFJAHAPM00000001",
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusBadRequest,
					Message: ErrUserIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
		{
			name:   "[003] Delete meal, mealId not indicated (400)",
			userID: "01FN3EEB2NVFJAHAPU00000001",
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusBadRequest,
					Message: ErrMealIDNotPresent.Error(),
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
		{
			name:   "[005] Meal does not exist (404)",
			userID: "01FN3EEB2NVFJAHAPU00000001",
			mealID: "01FN3EEB2NVFJAHAPM00000099",
			expectedResp: &ErrorResponse{
				Err: ErrorBody{
					Status:  http.StatusNotFound,
					Message: ErrMealNotFound.Error(),
				},
			},
			expectedStatusCode: http.StatusNotFound,
			wantErr:            true,
		},
	}
	getEchoContext := func(userId, mealId string) echo.Context {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, RouteMealID, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames(ParamUserID, ParamMealID)
		c.SetParamValues(userId, mealId)
		return c
	}
	for _, t := range tests {
		s.Run(t.name, func() {
			userManager := NewMealManager(*s.db)
			api := MealAPI{DB: *s.db, Manager: userManager}

			s.httpMock.On("GetUser", t.userID).Return(User{}, nil).Once()
			s.httpMock.On("GetCalendar", t.userID, Meal{Id: t.mealID}, true).Return(nil).Once()
			c := getEchoContext(t.userID, t.mealID)
			err := api.DeleteMealHandler(c)

			if t.wantErr {
				s.Equal(t.wantErr, err != nil)
				resp, ok := c.Response().Writer.(*httptest.ResponseRecorder)
				s.True(ok)
				body := resp.Body.Bytes()

				errorReturned := new(ErrorResponse)
				s.NoError(jsoniter.Unmarshal(body, errorReturned))
				s.Equal(errorReturned, t.expectedResp)
			}
			s.Equal(t.expectedStatusCode, c.Response().Status)
		})
	}
}
