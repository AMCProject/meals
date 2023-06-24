package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"meals/internal"
	"meals/internal/config"
	"meals/internal/handlers"
	"meals/internal/managers"
	"meals/pkg/database"
	"net/http"
)

const (
	banner = `
   ___    __  ___  _____        __  ___              __     
  / _ |  /  |/  / / ___/       /  |/  / ___  ___ _  / /  ___
 / __ | / /|_/ / / /__        / /|_/ / / -_)/ _  / / /  (_-<
/_/ |_|/_/  /_/  \___/       /_/  /_/  \__/ \_,_/ /_/  /___/

AMC Users Service
`
)

func main() {
	if err := config.LoadConfiguration(); err != nil {
		log.Fatal(err)
	}
	db := database.InitDB(config.Config.DBName)
	e := setUpServer(db)
	e.Logger.Fatal(e.Start(config.Config.Host + ":" + config.Config.Port))

}

func setUpServer(db *database.Database) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	addRoutes(e, *db)
	e.HideBanner = true
	fmt.Printf(banner)

	return e
}

func addRoutes(e *echo.Echo, db database.Database) {

	mealManager := managers.NewMealManager(db)
	externalManager := managers.NewExternalMealsManager()

	mealAPI := handlers.MealAPI{DB: db, Manager: mealManager, ExternalManager: externalManager}
	e.POST(internal.RouteMeal, mealAPI.PostMealHandler)
	e.GET(internal.RouteMeal, mealAPI.ListMealsHandler)
	e.GET(internal.RouteMealID, mealAPI.GetMealHandler)
	e.PUT(internal.RouteMealID, mealAPI.PutMealHandler)
	e.DELETE(internal.RouteMealID, mealAPI.DeleteMealHandler)

	e.GET(internal.RouteExternalMeals, mealAPI.GetExternalFoodsHandler)

	e.GET(internal.RouteIngredients, mealAPI.GetIngredients)
}
