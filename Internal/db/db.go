package db

import (
	"log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func DBconn() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("There was error connecting to the database: %v", err)
	  }
	  return db
}
