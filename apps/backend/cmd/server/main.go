package main

import (
	"circa/api"
	"circa/internal/config"
	"circa/internal/db"
	"circa/internal/handler"
	"circa/internal/service/auth"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log.Info().Msg("Starting Circa backend server...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading config")
	}

	// Initialize database
	_, err = db.InitPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing database")
	}

	// TODO: Initialize Redis when needed
	// redis.InitRedis(redis.RedisConfig{
	// 	Address: cfg.Redis.Address,
	// })

	// Initialize services
	authService := auth.NewService(15 * time.Minute) // Nonce expires in 15 minutes

	// Initialize handlers
	h := handler.NewHandler(authService)

	// Create Echo instance
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{cfg.FrontendURL},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderCookie},
		AllowCredentials: true,
	}))

	// Register OpenAPI handlers
	api.RegisterHandlers(e, h)

	// Start server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var g errgroup.Group

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
