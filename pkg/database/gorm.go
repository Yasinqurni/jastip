package database

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbMtx sync.Mutex
)

// Dialector for open connection based on driver
// Uncomment the code based on your database driver
// Run go get gorm.io/driver/{driver_name} to install the driver
// Example: go get gorm.io/driver/mysql and import "gorm.io/driver/mysql"
var dialectors = func(driver, dsn string) gorm.Dialector {
	return map[string]gorm.Dialector{
		"mysql": mysql.Open(dsn),
	}[driver]
}

// Connect and return *gorm.DB connection
func ConnectGorm(driver string, dbURL string) *gorm.DB {
	dialector := dialectors(driver, dbURL)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection error", err)
	}

	fmt.Println("Database connected successfully")
	return db
}

// Transaction will handle GORM database transaction
func Transaction(fn func(tx *gorm.DB) error, db *gorm.DB) error {
	dbMtx.Lock()
	defer dbMtx.Unlock()
	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true}).Begin()
	// If there are any panic it will rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Run the transaction
	err := fn(tx)
	if err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}
