package managers

import (
	"encoding/json"
	"meals/internal/models"
	"net/http"
)

const (
	APP_ID  = "e21c0369"
	APP_KEY = "165f842c901b3ba223f917cbe0a46ca9"
	ROUTE   = "https://api.edamam.com/api/recipes/v2?type=public"
)

var httpClient = &http.Client{}

type ExternalMealsManager struct {
}

type IExternalMealsManager interface {
	ListMeals(query string) ([]models.Meal, error)
}

func NewExternalMealsManager() *ExternalMealsManager {
	return &ExternalMealsManager{}
}
func (em *ExternalMealsManager) ListMeals(query string) ([]models.Meal, error) {
	url := ROUTE
	if query == "" {
		query = "pollo"
	}
	url += "&q=" + query + "&app_id=" + APP_ID + "&app_key=" + APP_KEY
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []models.Meal{}, err
	}
	var externalMeals models.ExternalMeals
	response, err := httpClient.Do(request)
	if err != nil {
		return []models.Meal{}, err
	}
	err = json.NewDecoder(response.Body).Decode(&externalMeals)
	if err != nil {
		return []models.Meal{}, err
	}
	var meals []models.Meal
	for _, m := range externalMeals.Hits {
		meals = append(meals, models.FromExternalToInternal(m.Recipe))
	}
	return meals, nil
}
