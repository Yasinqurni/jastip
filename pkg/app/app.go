package app

import (
	"jastip-app/config"

	"gorm.io/gorm"
)

// App DB
type Database struct {
	Gorm *gorm.DB
}

// App
type App struct {
	DB     Database
	Config *config.Config
}
