package handlers

import (
	"github.com/labstack/echo/v4"
	"meals/internal"
	"meals/internal/managers"
	"meals/internal/models"
	"meals/pkg/database"
	"meals/pkg/url"
	"net/http"
)

type MealAPI struct {
	DB              database.Database
	Manager         managers.IMealManager
	ExternalManager managers.IExternalMealsManager
}

func (a *MealAPI) GetMealHandler(c echo.Context) error {
	var userID, mealID string
	if err := url.ParseURLPath(c, url.PathMap{
		internal.ParamUserID: {Target: &userID, Err: internal.ErrUserIDNotPresent},
		internal.ParamMealID: {Target: &mealID, Err: internal.ErrMealIDNotPresent},
	}); err != nil {
		return internal.NewErrorResponse(c, err)
	}

	meal, err := a.Manager.GetMeal(userID, mealID)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}
	cleanMeal(meal)
	return c.JSON(http.StatusOK, meal)
}

func (a *MealAPI) ListMealsHandler(c echo.Context) error {
	var userID string

	if err := url.ParseURLPath(c, url.PathMap{
		internal.ParamUserID: {Target: &userID, Err: internal.ErrUserIDNotPresent},
	}); err != nil {
		return internal.NewErrorResponse(c, err)
	}

	filters := &models.MealsFilters{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, filters); err != nil {
		return internal.NewErrorResponse(c, err)
	}

	allMeals, err := a.Manager.ListMeals(userID, filters)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}
	if allMeals == nil {
		allMeals = []*models.Meal{}
	}
	for _, meal := range allMeals {
		cleanMeal(meal)
	}
	return c.JSON(http.StatusOK, allMeals)
}

func (a *MealAPI) PostMealHandler(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		internal.ParamUserID: {Target: &userID, Err: internal.ErrUserIDNotPresent},
	}); err != nil {
		return internal.NewErrorResponse(c, err)
	}

	mealFront := &models.Meal{}
	if err := c.Bind(mealFront); err != nil {
		return internal.NewErrorResponse(c, internal.ErrWrongBody)
	}
	meal, err := a.Manager.CreateMeal(userID, *mealFront)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}
	cleanMeal(meal)
	return c.JSON(http.StatusCreated, meal)

}

func (a *MealAPI) PutMealHandler(c echo.Context) error {
	var userID, mealID string
	if err := url.ParseURLPath(c, url.PathMap{
		internal.ParamUserID: {Target: &userID, Err: internal.ErrUserIDNotPresent},
		internal.ParamMealID: {Target: &mealID, Err: internal.ErrMealIDNotPresent},
	}); err != nil {
		return internal.NewErrorResponse(c, err)
	}

	mealFront := &models.Meal{}
	if err := c.Bind(mealFront); err != nil {
		return internal.NewErrorResponse(c, internal.ErrWrongBody)
	}
	mealFront.Id = mealID
	meal, err := a.Manager.UpdateMeal(userID, mealID, *mealFront)
	if err != nil {
		return internal.NewErrorResponse(c, err)

	}
	cleanMeal(meal)
	return c.JSON(http.StatusOK, meal)

}

func (a *MealAPI) DeleteMealHandler(c echo.Context) error {
	var userID, mealID string
	if err := url.ParseURLPath(c, url.PathMap{
		internal.ParamUserID: {Target: &userID, Err: internal.ErrUserIDNotPresent},
		internal.ParamMealID: {Target: &mealID, Err: internal.ErrMealIDNotPresent},
	}); err != nil {
		return internal.NewErrorResponse(c, err)
	}

	err := a.Manager.DeleteMeal(userID, mealID)
	if err != nil {
		return internal.NewErrorResponse(c, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (a *MealAPI) GetIngredients(c echo.Context) error {
	return c.JSON(http.StatusOK, models.Ingredients)
}

func (a *MealAPI) GetExternalFoodsHandler(c echo.Context) error {
	filters := &models.ExternalMealFilter{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, filters); err != nil {
		return internal.NewErrorResponse(c, err)
	}

	meals, err := a.ExternalManager.ListMeals(filters.Q)
	if err != nil {
		return internal.NewErrorResponse(c, internal.ErrorWithExternalAPI)
	}

	return c.JSON(http.StatusOK, meals)
}
func cleanMeal(meal *models.Meal) {
	meal.UserId = "" // Remove userID to prevent it from being serialized to JSON
}
