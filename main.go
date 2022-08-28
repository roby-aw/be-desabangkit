package main

import (
	"api-desatanggap/api"
	"api-desatanggap/app/modules"
	"api-desatanggap/config"
	"api-desatanggap/utils"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// @title API Poins ID
// @version 1.0
// @description Berikut API Poins ID
// @host api-poins-id.herokuapp.com/v1
// @BasePath /
func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	config := config.GetConfig()
	dbCon := utils.NewConnectionDatabase(config)

	defer dbCon.CloseConnection()

	controllers := modules.RegistrationModules(dbCon, config)

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339_nano}, method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "API Is Active")
	})
	api.RegistrationPath(e, controllers)
	// wg := sync.WaitGroup{}
	// wg.Add(1)
	go func() {
		if port == "" {
			port = "8080"
		}
		address := fmt.Sprintf(":%s", port)
		fmt.Println("goroutine jalan")

		if err := e.Start(address); err != nil {
			log.Fatal(err)
		}
	}()
	quit := make(chan os.Signal)
	<-quit
}
