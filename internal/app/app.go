// Package app configures and runs application.
package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"jastip-app/config"
	grpccontroller "jastip-app/internal/controller/grpc"
	v1 "jastip-app/internal/controller/http/v1"
	"jastip-app/internal/middleware"
	"jastip-app/internal/usecase"
	"jastip-app/pkg/app"
	"jastip-app/pkg/database"
	"jastip-app/pkg/httpserver"
	"jastip-app/pkg/logger"
	"jastip-app/pkg/upload"

	grpcserver "jastip-app/pkg/grpc/server"

	cors "github.com/rs/cors/wrapper/gin"
	"google.golang.org/grpc"
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

	// Uploader
	var uploader upload.Uploader
	uploader.ConnectS3(cfg)

	// App
	app := app.App{
		DB: app.Database{
			Gorm:  database.ConnectGorm("mysql", cfg.DB.URL),
			Mongo: database.ConnectMongo(context.Background(), cfg.Mongo.DSN, cfg.Mongo.DB),
		},
		Uploader: &uploader,
		Config:   cfg,
	}

	// Use case
	useCase := usecase.New(&app)

	// HTTP Server
	handler := gin.New()
	handler.Use(cors.AllowAll())
	handler.Use(middleware.GinLogger(ginLogger, errorLogger))
	v1.NewRouter(handler, useCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// gRPC Server
	grpcPort := fmt.Sprintf(":%s", cfg.GRPC.Port)
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal(err)
	}
	grpcSrv := grpc.NewServer()
	grpccontroller.NewGrpcServer(grpcSrv, useCase)
	grpcServer := grpcserver.New(
		grpcSrv,
		lis,
		grpcPort,
	)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println("app - Run - signal: ", s.String())
	case err = <-httpServer.Notify():
		log.Println(fmt.Errorf("app - Run - httpServer.Notify: %w", err).Error())
	case err = <-grpcServer.Notify():
		log.Println(fmt.Errorf("app - Run - grpcServer.Notify: %w", err).Error())
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err).Error())
	}
}
