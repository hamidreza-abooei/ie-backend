package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Id       uint   `gorm:"unique_index:index_addr_user"`
	Username string `gorm:"unique_index:index_addr_user"`
	Password string `gorm:"not null"`
	Urls     []Url  `gorm:"foreignkey:user_id"` // Personally I would resolve this one using DB query in order to database normalization
}
