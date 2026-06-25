package main

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/api/router"
	"Frank2006x/Pipe/internal/config"
	"Frank2006x/Pipe/internal/db"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := db.NewDB(config.POSTGRES_DB)
	if err != nil {
		panic(err)
	}
	queries := sqlc.New(db)
	router.RepositoryRouter(app, queries)
	app.Listen(":8080")

}
