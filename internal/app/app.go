// Package app configures and runs application.
package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"jastip-app/config"
	v1 "jastip-app/internal/controller/http/v1"
	"jastip-app/internal/middleware"
	"jastip-app/internal/usecase"
	"jastip-app/pkg/app"
	"jastip-app/pkg/database"
	"jastip-app/pkg/httpserver"
	"jastip-app/pkg/logger"

	cors "github.com/rs/cors/wrapper/gin"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	errorLogger := logger.New(cfg.Log.Level)
	ginLogger := logger.New("request")
	logger.NewGlobal(errorLogger)

	// Set timezone
	os.Setenv("TZ", cfg.App.TimeZone)
	krTime, err := time.LoadLocation(cfg.App.TimeZone)
	if err != nil {
		log.Fatal("Timezone error: ", err)
	}
	time.Local = krTime

	// App
	app := app.App{
		DB: app.Database{
			Gorm: database.ConnectGorm("mysql", cfg.DB.URL),
		},
		Config: cfg,
	}

	// Use case
	useCase := usecase.New(&app)

	// HTTP Server
	handler := gin.New()
	handler.Use(cors.AllowAll())
	handler.Use(middleware.GinLogger(ginLogger, errorLogger))
	v1.NewRouter(handler, useCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: ", s.String())
	case err = <-httpServer.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	}
}
