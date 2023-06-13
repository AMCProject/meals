package internal

import (
	"github.com/go-playground/validator/v10"
	"meals/pkg/database"
)

type MealManager struct {
	db             *SQLiteMealRepository
	validate       *validator.Validate
	allIngredients map[string]int
}

var Microservices EndpointsI = &Endpoints{}

type IMealManager interface {
	GetMeal(userID, mealID string) (meal *Meal, err error)
	ListMeals(userID string, filters *MealsFilters) (meals []*Meal, err error)
	UpdateMeal(userID string, mealID string, mealPut Meal) (meal *Meal, err error)
	CreateMeal(userID string, mealPost Meal) (meal *Meal, err error)
	DeleteMeal(userID, mealID string) (err error)
}

func NewMealManager(db database.Database) *MealManager {
	return &MealManager{
		db:             NewSQLiteMealRepository(&db),
		validate:       validator.New(),
		allIngredients: GetAllIngredients(),
	}
}

// GetMeal function to get a specific meal from a user
func (m *MealManager) GetMeal(userID, mealID string) (meal *Meal, err error) {
	if _, err = Microservices.GetUser(userID); err != nil {
		return nil, err
	}
	return m.db.GetMeal(userID, mealID)
}

// ListMeals returns all the meals created by a user
func (m *MealManager) ListMeals(userID string, filters *MealsFilters) (meals []*Meal, err error) {
	if _, err = Microservices.GetUser(userID); err != nil {
		return nil, err
	}
	return m.db.ListMeals(userID, *filters)
}

// UpdateMeal function to update the meal selected (if any parameter is missing we get the oldest ones
func (m *MealManager) UpdateMeal(userID string, mealID string, mealPut Meal) (meal *Meal, err error) {
	if _, err = Microservices.GetUser(userID); err != nil {
		return nil, err
	}

	if err = m.validate.Struct(mealPut); err != nil {
		return nil, err
	}
	mealGet, err := m.db.GetMeal(userID, mealID)
	if err != nil {
		return nil, err
	}

	var kcal int
	for _, ing := range mealPut.Ingredients {
		kcal += m.allIngredients[ing]
	}
	mealPut.Kcal = kcal / len(mealPut.Ingredients)

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
func (m *MealManager) CreateMeal(userID string, mealPost Meal) (meal *Meal, err error) {
	if _, err = Microservices.GetUser(userID); err != nil {
		return nil, err
	}

	if err = m.validate.Struct(mealPost); err != nil {
		return nil, ErrWrongBody
	}
	var kcal int
	for _, ing := range mealPost.Ingredients {
		kcal += m.allIngredients[ing]
	}
	mealPost.Kcal = kcal / len(mealPost.Ingredients)

	return m.db.CreateMeal(userID, mealPost)
}

// DeleteMeal function to delete on meal from user
func (m *MealManager) DeleteMeal(userID, mealID string) (err error) {
	if _, err = Microservices.GetUser(userID); err != nil {
		return err
	}
	if _, err = m.db.GetMeal(userID, mealID); err != nil {
		return err
	}
	if err = m.db.DeleteMeal(userID, mealID); err != nil {
		return err
	}
	if err = Microservices.GetCalendar(userID, Meal{Id: mealID}, true); err != nil {
		return err
	}
	return
}

func GetAllIngredients() map[string]int {
	allIngredients := make(map[string]int)
	for _, ing := range Ingredients {
		for value, key := range ing {
			allIngredients[value] = key
		}
	}
	return allIngredients
}
