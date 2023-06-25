package managers

import (
	"github.com/go-playground/validator/v10"
	"meals/internal"
	"meals/internal/models"
	"meals/internal/repositories"
	"meals/internal/utils"
	"meals/pkg/database"
	"reflect"
)

type MealManager struct {
	db             *repositories.SQLiteMealRepository
	validate       *validator.Validate
	allIngredients map[string]int
}

var Microservices utils.EndpointsI = &utils.Endpoints{}

type IMealManager interface {
	GetMeal(userID, mealID string) (meal *models.Meal, err error)
	ListMeals(userID string, filters *models.MealsFilters) (meals []*models.Meal, err error)
	UpdateMeal(userID string, mealID string, mealPut models.Meal) (meal *models.Meal, err error)
	CreateMeal(userID string, mealPost models.Meal) (meal *models.Meal, err error)
	DeleteMeal(userID, mealID string) (err error)
}

func NewMealManager(db database.Database) *MealManager {
	return &MealManager{
		db:             repositories.NewSQLiteMealRepository(&db),
		validate:       validator.New(),
		allIngredients: GetAllIngredients(),
	}
}

// GetMeal function to get a specific meal from a user
func (m *MealManager) GetMeal(userID, mealID string) (meal *models.Meal, err error) {
	return m.db.GetMeal(userID, mealID)
}

// ListMeals returns all the meals created by a user
func (m *MealManager) ListMeals(userID string, filters *models.MealsFilters) (meals []*models.Meal, err error) {
	return m.db.ListMeals(userID, *filters)
}

// UpdateMeal function to update the meal selected (if any parameter is missing we get the oldest ones
func (m *MealManager) UpdateMeal(userID string, mealID string, mealPut models.Meal) (meal *models.Meal, err error) {

	if err = m.validate.Struct(mealPut); err != nil {
		return nil, err
	}
	mealGet, err := m.db.GetMeal(userID, mealID)
	if err != nil {
		return nil, err
	}

	var kcal int
	if !reflect.DeepEqual(mealPut.Ingredients, mealGet.Ingredients) {
		for _, ing := range mealPut.Ingredients {
			kcal += m.allIngredients[ing]
		}
		mealPut.Kcal = kcal / len(mealPut.Ingredients)
	}

	meal, err = m.db.UpdateMeal(userID, mealID, mealPut)
	if err != nil {
		return nil, err
	}
	if meal.Name != mealGet.Name {
		if err = Microservices.GetCalendar(userID, *meal, false); err != nil {
			return nil, err
		}
	}
	return

}

// CreateMeal function to create a new meal for the user selected
func (m *MealManager) CreateMeal(userID string, mealPost models.Meal) (meal *models.Meal, err error) {

	if err = m.validate.Struct(mealPost); err != nil {
		return nil, internal.ErrWrongBody
	}
	_, err = m.db.GetMealByName(userID, mealPost.Name)
	if err != nil {
		return nil, err
	}
	var kcal int
	if mealPost.Kcal == 0 {
		for _, ing := range mealPost.Ingredients {
			kcal += m.allIngredients[ing]
		}
		mealPost.Kcal = kcal / len(mealPost.Ingredients)
	}
	return m.db.CreateMeal(userID, mealPost)
}

// DeleteMeal function to delete on meal from user
func (m *MealManager) DeleteMeal(userID, mealID string) (err error) {
	if _, err = m.db.GetMeal(userID, mealID); err != nil {
		return err
	}
	if err = m.db.DeleteMeal(userID, mealID); err != nil {
		return err
	}
	if err = Microservices.GetCalendar(userID, models.Meal{Id: mealID}, true); err != nil {
		return err
	}

	return
}

func GetAllIngredients() map[string]int {
	allIngredients := make(map[string]int)
	for _, ing := range models.Ingredients {
		for value, key := range ing {
			allIngredients[value] = key
		}
	}
	return allIngredients
}
