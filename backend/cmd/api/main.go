package main

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/api/middleware"
	"Frank2006x/Pipe/internal/api/router"
	"Frank2006x/Pipe/internal/config"
	"Frank2006x/Pipe/internal/db"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())
	app.Use(middleware.RequestIDMiddleware())
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	err = config.CheckConfig(cfg)
	if err != nil {
		panic(err)
	}

	db, err := db.NewDB(cfg.POSTGRES_DB)
	if err != nil {
		panic(err)
	}
	queries := sqlc.New(db)
	router.RepositoryRouter(app, queries)
	router.AuthRouter(app, queries, cfg)
	app.Listen(":8080")

}
