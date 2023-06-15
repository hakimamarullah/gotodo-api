package models

import "gorm.io/gorm"

type TodoItem struct {
	gorm.Model
	Description string
	Completed   bool `gorm:"default=false"`
	UserID      uint `gorm:"foreignKey:UserID;references:ID" json:"userId"`
}
