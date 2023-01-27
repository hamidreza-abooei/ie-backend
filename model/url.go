package model

import "github.com/jinzhu/gorm"

type Url struct {
	gorm.Model
	Id          uint   `gorm:"unique_index:index_addr_user"`
	UserId      uint   `gorm:"unique_index:index_addr_user"`
	Adress      string `gorm:"unique_index:index_addr_user"`
	Threshold   int
	FailedTimes int
	Requests    []Request `gorm:"foreignkey:url_id"` // Personally I would resolve this one using DB query in order to database normalization
}
