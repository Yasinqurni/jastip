package app

import (
	"jastip-app/config"
	"jastip-app/pkg/upload"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// App DB
type Database struct {
	Gorm  *gorm.DB
	Mongo *mongo.Database
}

// App
type App struct {
	DB       Database
	Uploader *upload.Uploader
	Config   *config.Config
}
