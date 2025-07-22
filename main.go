package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lab1702/goduck/internal/config"
	"github.com/lab1702/goduck/internal/database"
	"github.com/lab1702/goduck/internal/handlers"
	"github.com/lab1702/goduck/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logrus.WithError(err).Fatal("Configuration validation failed")
	}

	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(level)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})

	db, err := database.NewDB(cfg.DatabasePath, cfg.MaxConnections, cfg.ReadWrite)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize database")
	}
	defer db.Close()

	queryHandler := handlers.NewQueryHandler(db, cfg.QueryTimeout)

	// Rate limiter: 60 requests per minute per IP
	rateLimiter := middleware.NewRateLimiter(60)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(rateLimiter.Middleware())

	router.POST("/query", queryHandler.ExecuteQuery)
	router.GET("/health", queryHandler.Health)
	router.GET("/metrics", queryHandler.Metrics)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logrus.WithField("port", cfg.Port).Info("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.WithError(err).Error("Server forced to shutdown")
	} else {
		logrus.Info("Server shutdown complete")
	}
}
