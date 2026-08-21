package main

import (
	"Frank2006x/Pipe/db/sqlc"
	"Frank2006x/Pipe/internal/auth"
	"Frank2006x/Pipe/internal/cache"
	"Frank2006x/Pipe/internal/config"
	"Frank2006x/Pipe/internal/db"
	"Frank2006x/Pipe/internal/executor"
	"Frank2006x/Pipe/internal/github"
	"Frank2006x/Pipe/internal/queue"
	"Frank2006x/Pipe/internal/server/middleware"
	"Frank2006x/Pipe/internal/server/router"
	"Frank2006x/Pipe/internal/service"
	"Frank2006x/Pipe/internal/worker"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
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
	rabbitmq, err := queue.NewRabbitmq(cfg.RABBITMQ_URL)
	if err != nil {
		panic(err)
	}
	defer rabbitmq.Close()

	redisCache, err := cache.New(cfg.REDIS_URL)
	if err != nil {
		log.Printf("[WARNING] Failed to connect to Redis, continuing without cache: %v", err)
	}
	if redisCache != nil {
		defer redisCache.Close()
	}

	queries := sqlc.New(pool)
	githubClient := github.NewClient(cfg)
	jwtMaker := auth.NewJwtMaker(cfg.JWT_SECRET)
	githubService := service.NewGithubService(githubClient, queries, &cfg)
	pipelineService := service.NewPipelineService(queries, pool, rabbitmq)
	appCtx, cancel := context.WithCancel(context.Background())
	dockerExecutor, err := executor.NewDockerExecutor()
	if err != nil {
		log.Printf("[WARNING] Failed to initialize Docker executor: %v", err)
	}

	worker := worker.NewPipelineWorker(rabbitmq, pipelineService, dockerExecutor)

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
	}))

	var workerWg sync.WaitGroup
	workerWg.Add(1)
	go func() {
		defer workerWg.Done()
		if err := worker.Start(appCtx); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("[ERROR] Worker STOP completely: %v", err)
		}
	}()
	app.Use(logger.New())
	app.Use(middleware.RequestIDMiddleware())

	router.RepositoryRouter(app, queries, githubService, jwtMaker, pool)
	router.AuthRouter(app, queries, githubClient, jwtMaker, redisCache)
	router.GithubRouter(app, githubService, jwtMaker)
	router.WebhookRouter(app, queries, pipelineService)
	router.PipelineRouter(app, pipelineService, jwtMaker)
	app.Get("/ping", func(c fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})


	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		log.Println("[INFO] Fiber web server starting up on port :8080...")
		if err := app.Listen(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("[FATAL] Web server crashed completely: %v", err)
		}
	}()

	<-shutdownChannel
	log.Println("[INFO] Fiber web server shutting down...")
	cancel() // Notify worker to stop consuming and executing

	log.Println("[INFO] Waiting for worker to finish execution...")
	workerWg.Wait()
	log.Println("[SUCCESS] Worker shutdown complete.")

	shutDownCtx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelTimeout()

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
