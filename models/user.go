package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex"`
	Password string
	TodoItem []TodoItem `gorm:"constraint:OnDelete:CASCADE"`
}
