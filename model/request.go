package model

import "github.com/jinzhu/gorm"

type Request struct {
	gorm.Model
	Id     uint `gorm:"unique_index:index_addr_user"` // for preventing url duplication for a single user
	UrlId  uint
	Result int
}
