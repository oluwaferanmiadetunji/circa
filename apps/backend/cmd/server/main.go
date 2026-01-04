package main

import (
	"circa/api"
	"circa/internal/config"
	"circa/internal/db"
	"circa/internal/email"
	"circa/internal/handler"
	"circa/internal/queue"
	"circa/internal/redis"
	"circa/internal/service/auth"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func main() {

	log.Info().Msg("Starting Circa backend server...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}

	// Initialize database
	store, err := db.InitPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}

	// Initialize Redis (optional - app can run without it)
	if err := redis.InitRedis(redis.RedisConfig{
		Address:  cfg.Redis.Address,
		Password: cfg.Redis.Password,
	}); err != nil {
		log.Warn().Err(err).Msg("Failed to connect to Redis, continuing without Redis support")
	} else {
		log.Info().Msg("Redis connected successfully")
	}

	queueService := queue.NewService(store)
	emailService := email.NewService(cfg.ResendAPIKey)
	queueWorker := queue.NewWorker(queueService, emailService)

	authService := auth.NewService(store, queueService, cfg.FrontendURL, 15*time.Minute)

	// Initialize handlers
	h := handler.NewHandler(authService, cfg)

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{cfg.FrontendURL},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderCookie},
		AllowCredentials: true,
	}))

	// Register OpenAPI handlers
	api.RegisterHandlers(e, h)

	// Start server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var g errgroup.Group

	g.Go(func() error {
		queueWorker.Start(ctx)
		return nil
	})

	// Start HTTP server
	g.Go(func() error {
		serverAddr := ":" + cfg.Port
		log.Info().Str("address", serverAddr).Msg("Starting HTTP server")
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("HTTP server error")
			return err
		}
		return nil
	})

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	g.Go(func() error {
		select {
		case <-sigChan:
			log.Info().Msg("Received termination signal, shutting down...")
			cancel()
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	// Graceful shutdown
	g.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := e.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("Error during server shutdown")
			return err
		}

		log.Info().Msg("Server shutdown complete")
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatal().Err(err).Msg("Server shut down with error")
	}

	log.Info().Msg("Server gracefully stopped")
}
