package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Scorzoner/effective-mobile-test/internal/api/handlers"
	"github.com/Scorzoner/effective-mobile-test/internal/api/router"
	"github.com/Scorzoner/effective-mobile-test/internal/config"
	"github.com/Scorzoner/effective-mobile-test/internal/database"
	"github.com/Scorzoner/effective-mobile-test/internal/logger"
	"github.com/golang-migrate/migrate/v4"
)

func Run() {
	// start logger
	err := logger.Init()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize logger: %w", err))
	}

	logger.Zap.Info("Initialized logger")

	// load config
	logger.Zap.Info("Loading config")
	cfg, err := config.Load()
	if err != nil {
		logger.Zap.Fatal(fmt.Errorf("failed to load config: %w", err))
	}
	logger.Zap.Info("Config loaded: ", fmt.Sprintf("%+v", cfg))

	// open db connection
	logger.Zap.Info("Opening database connection")
	db, err := database.Open(cfg)
	if err != nil {
		logger.Zap.Fatal(fmt.Errorf("failed to open pgsql connection: %w", err))
	}
	defer db.Close()

	// run migrations
	logger.Zap.Info("Running migrations")
	err = database.RunMigrations(db)
	if err != nil && err != migrate.ErrNoChange {
		logger.Zap.Fatal(fmt.Errorf("failed to run migrations: %w", err))
	}

	// initialize router/handlers
	logger.Zap.Info("Initializing handlers")
	hq, err := handlers.NewHandlerQueries(db, cfg)
	if err != nil {
		logger.Zap.Fatal(fmt.Errorf("failed to initialize queries: %w", err))
	}

	r := router.New(hq)

	// start server
	logger.Zap.Info("Configuring and starting the server")
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	// graceful shutdown
	shutdownError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		signal := <-quit

		logger.Zap.Info(fmt.Sprintf("Shutting down server gracefully: %s", signal.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	logger.Zap.Info(fmt.Sprintf("Server is running on port: %d", cfg.Port))

	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Zap.Info(fmt.Errorf("HTTP server error: %w", err))
		return
	}

	errShut := <-shutdownError
	if errShut != nil {
		logger.Zap.Fatal(fmt.Errorf("HTTP shutdown error: %w", errShut))
	}
	logger.Zap.Info("Graceful shutdown complete")
}
