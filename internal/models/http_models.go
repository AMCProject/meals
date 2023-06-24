package models

type User struct {
	Id       string  `db:"id" json:"id,omitempty"`
	Name     *string `db:"name" json:"name,omitempty"`
	Mail     string  `db:"mail" json:"mail" validate:"required,excludes= "`
	Password string  `db:"password" json:"password" validate:"required,excludes= "`
}

type Calendar struct {
	UserId string `db:"user_id" json:"user_id"`
	MealId string `db:"meal_id" json:"meal_id"`
	Name   string `json:"name" json:"name"`
	Date   string `db:"date" json:"date"`
}

type UpdateWeekCalendar struct {
	From string `json:"from"`
	To   string `json:"to"`
}
