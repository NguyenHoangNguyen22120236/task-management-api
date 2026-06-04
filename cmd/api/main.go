package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	middleware "github.com/oapi-codegen/gin-middleware"

	"github.com/NguyenNH36/task-management-api/internal/api"
	"github.com/NguyenNH36/task-management-api/internal/config"
	"github.com/NguyenNH36/task-management-api/internal/database"
	"github.com/NguyenNH36/task-management-api/internal/handler"
	"github.com/NguyenNH36/task-management-api/internal/logger"
	customMiddleware "github.com/NguyenNH36/task-management-api/internal/middleware"
	"github.com/NguyenNH36/task-management-api/internal/repository"
	"github.com/NguyenNH36/task-management-api/internal/service"
	"github.com/getkin/kin-openapi/openapi3"
)

func main() {
	configuration := config.Load()

	gin.SetMode(configuration.GinMode)

	appLogger := logger.New()

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(customMiddleware.RequestID())
	router.Use(customMiddleware.RequestLogger(appLogger))

	swagger, err := openapi3.NewLoader().LoadFromFile("api/openapi.yaml")
	if err != nil {
		log.Fatalf("failed to load OpenAPI spec: %v", err)
	}

	swagger.Servers = nil

	router.Use(middleware.OapiRequestValidator(swagger))

	db, err := database.NewPostgresPool(context.Background(), configuration.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	taskRepo := repository.NewPostgresTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	api.RegisterHandlers(router, taskHandler)

	server := &http.Server{
		Addr:              ":" + configuration.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		appLogger.Info("server_started", "port", configuration.Port)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appLogger.Error("server_failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	receivedSignal := <-quit

	appLogger.Info("shutdown_signal_received", "signal", receivedSignal.String())

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("server_forced_to_shutdown", "error", err)
		os.Exit(1)
	}

	appLogger.Info("server_shutdown_completed")
}