package data

import (
	"database/sql"
	"gorm.io/gorm"
	"log"
)

func getDB(db *gorm.DB) *sql.DB {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("failed to get sql.DB from gorm.DB", err)
	}
	return sqlDB
}
