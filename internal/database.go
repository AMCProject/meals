package internal

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"meals/pkg/database"
	"time"
)

const (
	getMeal    = "SELECT * FROM meals WHERE user_id = ? AND id = ?"
	listMeals  = "SELECT * FROM meals WHERE user_id = ? "
	updateMeal = "UPDATE meals SET name = ?, description = ?, image = ?, type = ?, ingredients = ?, kcal = ?, seasons = ? WHERE user_id = ? AND id = ?"
	createMeal = "INSERT INTO meals(id,user_id,name,description,image,type,ingredients,kcal,seasons) VALUES (?,?,?,?,?,?,?,?,?)"
	deleteMeal = "DELETE FROM meals WHERE user_id = ? AND id = ?"
)

type MealRepository interface {
	GetMeal(userID, mealID string) (meal *Meal, err error)
	ListMeals(userID string, filters MealsFilters) (meals []*Meal, err error)
	UpdateMeal(userID string, mealID string, mealPut Meal) (meal *Meal, err error)
	CreateMeal(userID string, mealPost Meal) (meal *Meal, err error)
	DeleteMeal(userID, mealID string) (err error)
}

type SQLiteMealRepository struct {
	db *database.Database
}

func NewSQLiteMealRepository(db *database.Database) *SQLiteMealRepository {
	return &SQLiteMealRepository{
		db: db,
	}
}

func (r *SQLiteMealRepository) GetMeal(userId, mealId string) (*Meal, error) {
	var mealsAux []MealDB
	err := r.db.Conn.Select(&mealsAux, getMeal, userId, mealId)
	if err != nil {
		log.Error(err)
		return nil, ErrSomethingWentWrong
	}
	if len(mealsAux) == 0 {
		return nil, ErrMealNotFound
	}
	return MealToAPI(&mealsAux[0]), nil

}

func (r *SQLiteMealRepository) ListMeals(userId string, filters MealsFilters) (meals []*Meal, err error) {
	var mealsDB []MealDB
	err = r.db.Conn.Select(&mealsDB, listMeals+applyFilters(filters), userId)
	if err != nil {
		log.Error(err)
		return nil, ErrSomethingWentWrong
	}
	if len(mealsDB) == 0 {
		return nil, ErrMealNotFound
	}
	for _, m := range mealsDB {
		meals = append(meals, MealToAPI(&m))
	}
	return
}

func (r *SQLiteMealRepository) CreateMeal(userID string, mealPost Meal) (*Meal, error) {
	id, _ := ulid.New(ulid.Now(), ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0))
	mealPost.Id = id.String()
	mealDB := MealFromAPI(&mealPost)
	_, err := r.db.Conn.Exec(createMeal, mealDB.Id, userID, mealDB.Name, mealDB.Description, mealDB.Image, mealDB.Type, mealDB.Ingredients, mealDB.Kcal, mealDB.Seasons)
	if err != nil {
		log.Error(err)
		return nil, ErrSomethingWentWrong
	}

	return &mealPost, nil
}

func (r *SQLiteMealRepository) UpdateMeal(userID string, mealID string, mealUpdate Meal) (meal *Meal, err error) {
	mealDB := MealFromAPI(&mealUpdate)
	_, err = r.db.Conn.Exec(updateMeal, mealDB.Name, mealDB.Description, mealDB.Image, mealDB.Type, mealDB.Ingredients, mealDB.Kcal, mealDB.Seasons, userID, mealID)
	if err != nil {
		log.Error(err)
		return nil, ErrSomethingWentWrong
	}
	return &mealUpdate, nil
}

func (r *SQLiteMealRepository) DeleteMeal(userID, mealID string) (err error) {
	_, err = r.db.Conn.Exec(deleteMeal, userID, mealID)
	if err != nil {
		log.Error(err)
		return ErrSomethingWentWrong
	}
	return
}

func applyFilters(filters MealsFilters) (query string) {
	if filters.Name != nil {
		*filters.Name = "%" + *filters.Name + "%"
		query += fmt.Sprintf("AND name LIKE '%s'", *filters.Name)
	}
	if filters.Type != nil {
		query += fmt.Sprintf("AND type = '%s'", *filters.Type)
	}
	if filters.Healthy != nil && *filters.Healthy {
		query += "ORDER BY kcal ASC"
	}
	if len(filters.Season) > 0 {
		for _, s := range filters.Season {
			season := "%" + s + "%"
			query += fmt.Sprintf("AND seasons LIKE '%s'", season)
		}
		general := "%general%"
		query += fmt.Sprintf("OR seasons LIKE '%s'", general)
	}
	return
}
