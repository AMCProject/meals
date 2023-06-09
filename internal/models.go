package internal

import "strings"

var Ingredients = map[string]int{
	"Garbanzos":      361,
	"Judías":         343,
	"Lentejas":       336,
	"Cebolla":        47,
	"Calabaza":       24,
	"Calabacín":      31,
	"Champiñones":    24,
	"Coliflor":       30,
	"Espárragos":     26,
	"Espinacas":      32,
	"Lechuga":        18,
	"Pepino":         12,
	"Tomate":         22,
	"Zanahoria":      42,
	"Nata":           298,
	"Queso blanco":   70,
	"Queso":          300,
	"Bacon":          665,
	"Cerdo":          220,
	"Pollo":          145,
	"Ternera":        176,
	"Atún":           173,
	"Calamar":        82,
	"Salmón":         165,
	"Crustáceo":      79,
	"Pescado blanco": 211,
	"Arroz":          352,
	"Pasta":          368,
	"Huevo duro":     147,
	"Huevo frito":    340,
	"Patata frita":   312,
	"Patata cocida":  85,
	"Patata asada":   90}

type MealDB struct {
	Id          string `db:"id" json:"id,omitempty"`
	UserId      string `db:"user_id" json:"userId"`
	Name        string `db:"name" json:"name" validate:"required"`
	Description string `db:"description" json:"description,omitempty"`
	Image       string `db:"image" json:"image,omitempty"`
	Type        string `db:"type" json:"type" validate:"required,oneof=weekly occasional standard"`
	Ingredients string `db:"ingredients" json:"ingredients" validate:"required"`
	Kcal        int    `db:"kcal" json:"kcal"`
	Seasons     string `db:"seasons" json:"seasons"`
	//Creator     int    `db:"creator" json:"creator"`
	//Saves       int    `db:"saves" json:"saves"`
}

type Meal struct {
	Id          string   `json:"id"`
	UserId      string   `json:"userId"`
	Name        string   `json:"name" validate:"required"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Type        string   `json:"type" validate:"required,oneof=weekly occasional standard"`
	Ingredients []string `json:"ingredients" validate:"required"`
	Kcal        int      `json:"kcal"`
	Seasons     []string `json:"seasons" validate:"required,dive,oneof=summer winter spring fall general"`
	//Creator     int      `json:"creator"`
	//Saves       int      `json:"saves"`
}

type MealsFilters struct {
	Name    *string  `query:"name"`
	Type    *string  `query:"type"`
	Healthy *bool    `query:"healthy"`
	Season  []string `query:"[]season"`
}

func MealToAPI(meal *MealDB) *Meal {
	return &Meal{
		Id:          meal.Id,
		UserId:      meal.UserId,
		Name:        meal.Name,
		Description: meal.Description,
		Image:       meal.Image,
		Type:        meal.Type,
		Ingredients: strings.Split(meal.Ingredients, ","),
		Kcal:        meal.Kcal,
		Seasons:     strings.Split(meal.Seasons, ","),
	}
}

func MealFromAPI(meal *Meal) *MealDB {
	return &MealDB{
		Id:          meal.Id,
		UserId:      meal.UserId,
		Name:        meal.Name,
		Description: meal.Description,
		Image:       meal.Image,
		Type:        meal.Type,
		Ingredients: strings.Join(meal.Ingredients, ","),
		Kcal:        meal.Kcal,
		Seasons:     strings.Join(meal.Seasons, ","),
	}
}

//definitions for endpoint calls//

type User struct {
	Id       string  `db:"id" json:"id,omitempty"`
	Name     *string `db:"name" json:"name,omitempty"`
	Mail     string  `db:"mail" json:"mail" validate:"required,excludes= "`
	Password string  `db:"password" json:"password" validate:"required,excludes= "`
}
