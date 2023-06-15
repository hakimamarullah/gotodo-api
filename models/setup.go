package models

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectPostgreSQL() {
	dsn := "host=localhost user=go password=go dbname=gotodo port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Error(err)
	}

	db.AutoMigrate(&TodoItem{})

	DB = db
}
