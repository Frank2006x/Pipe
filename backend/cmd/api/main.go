package main

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/api/middleware"
	"Frank2006x/Pipe/internal/api/router"
	"Frank2006x/Pipe/internal/auth"
	"Frank2006x/Pipe/internal/config"
	"Frank2006x/Pipe/internal/db"
	"Frank2006x/Pipe/internal/github"
	"Frank2006x/Pipe/internal/service"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	err = config.CheckConfig(cfg)
	if err != nil {
		panic(err)
	}
	pool, err := db.NewDB(cfg.POSTGRES_DB)
	if err != nil {
		panic(err)
	}
	queries := sqlc.New(pool)
	githubClient := github.NewClient(cfg)
	jwtMaker := auth.NewJwtMaker(cfg.JWT_SECRET)
	githubService := service.NewGithubService(githubClient, queries)
	app := fiber.New()
	app.Use(logger.New())
	app.Use(middleware.RequestIDMiddleware())

	router.RepositoryRouter(app, queries,githubService, jwtMaker)
	router.AuthRouter(app, queries, githubClient, jwtMaker)
	router.GithubRouter(app, githubService, jwtMaker)

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Println("[INFO] Fiber web server starting up on port :8080...")
		if err := app.Listen(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[FATAL] Web server crashed completely: %v", err)
		}
	}()

	<-shutdownChannel
	log.Println("[INFO] Fiber web server shutting down...")

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("[INFO] Step 1/2: Draining and shutting down incoming HTTP traffic pools...")
	if err := app.ShutdownWithContext(shutDownCtx); err != nil {
		log.Printf("[ERROR] Fiber failed to shut down perfectly within deadline: %v", err)
	} else {
		log.Println("[SUCCESS] HTTP server safely deactivated.")
	}

	log.Println("[INFO] Step 2/2: Severing open connections to the PostgreSQL database...")
	pool.Close()
	log.Println("[SUCCESS] Postgres connection array successfully closed down.")

	log.Println("[INFO] All core dependencies decommissioned cleanly. App exit safe.")

}
