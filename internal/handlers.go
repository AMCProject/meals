package internal

import (
	"github.com/labstack/echo/v4"
	"meals/pkg/database"
	"meals/pkg/url"
	"net/http"
)

type MealAPI struct {
	DB      database.Database
	Manager IMealManager
}

func (a *MealAPI) GetMealHandler(c echo.Context) error {
	var userID, mealID string
	if err := url.ParseURLPath(c, url.PathMap{
		ParamUserID: {Target: &userID, Err: ErrUserIDNotPresent},
		ParamMealID: {Target: &mealID, Err: ErrMealIDNotPresent},
	}); err != nil {
		return NewErrorResponse(c, err)
	}

	meal, err := a.Manager.GetMeal(userID, mealID)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	cleanMeal(meal)
	return c.JSON(http.StatusOK, meal)
}

func (a *MealAPI) ListMealsHandler(c echo.Context) error {
	var userID string

	if err := url.ParseURLPath(c, url.PathMap{
		ParamUserID: {Target: &userID, Err: ErrUserIDNotPresent},
	}); err != nil {
		return NewErrorResponse(c, err)
	}

	filters := &MealsFilters{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, filters); err != nil {
		return NewErrorResponse(c, err)
	}

	allMeals, err := a.Manager.ListMeals(userID, filters)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	if allMeals == nil {
		allMeals = []*Meal{}
	}
	for _, meal := range allMeals {
		cleanMeal(meal)
	}
	return c.JSON(http.StatusOK, allMeals)
}

func (a *MealAPI) PostMealHandler(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		ParamUserID: {Target: &userID, Err: ErrUserIDNotPresent},
	}); err != nil {
		return NewErrorResponse(c, err)
	}

	mealFront := &Meal{}
	if err := c.Bind(mealFront); err != nil {
		return NewErrorResponse(c, ErrWrongBody)
	}
	meal, err := a.Manager.CreateMeal(userID, *mealFront)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	cleanMeal(meal)
	return c.JSON(http.StatusCreated, meal)

}

func (a *MealAPI) PutMealHandler(c echo.Context) error {
	var userID, mealID string
	if err := url.ParseURLPath(c, url.PathMap{
		ParamUserID: {Target: &userID, Err: ErrUserIDNotPresent},
		ParamMealID: {Target: &mealID, Err: ErrMealIDNotPresent},
	}); err != nil {
		return NewErrorResponse(c, err)
	}

	mealFront := &Meal{}
	if err := c.Bind(mealFront); err != nil {
		return NewErrorResponse(c, ErrWrongBody)
	}
	mealFront.Id = mealID
	meal, err := a.Manager.UpdateMeal(userID, mealID, *mealFront)
	if err != nil {
		return NewErrorResponse(c, err)

	}
	cleanMeal(meal)
	return c.JSON(http.StatusOK, meal)

}

func (a *MealAPI) DeleteMealHandler(c echo.Context) error {
	var userID, mealID string
	if err := url.ParseURLPath(c, url.PathMap{
		ParamUserID: {Target: &userID, Err: ErrUserIDNotPresent},
		ParamMealID: {Target: &mealID, Err: ErrMealIDNotPresent},
	}); err != nil {
		return NewErrorResponse(c, err)
	}

	err := a.Manager.DeleteMeal(userID, mealID)
	if err != nil {
		return NewErrorResponse(c, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (a *MealAPI) GetIngredients(c echo.Context) error {
	return c.JSON(http.StatusOK, Ingredients)
}

func (a *MealAPI) GetOtherMeals(c echo.Context) error {
	var userID string
	if err := url.ParseURLPath(c, url.PathMap{
		ParamUserID: {Target: &userID, Err: ErrUserIDNotPresent},
	}); err != nil {
		return NewErrorResponse(c, err)
	}

	allMeals, err := a.Manager.GetOtherMeals(userID)
	if err != nil {
		return NewErrorResponse(c, err)
	}
	return c.JSON(http.StatusOK, allMeals)

}

func cleanMeal(meal *Meal) {
	meal.UserId = "" // Remove userID to prevent it from being serialized to JSON
}
