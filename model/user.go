package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Id       uint   `gorm:"unique_index:not_null"`
	Username string `gorm:"unique_index:not_null"`
	Password string `gorm:"not null"`
	Urls     []Url  `gorm:"foreignkey:user_id"` // Personally I would resolve this one using DB query in order to database normalization
}
