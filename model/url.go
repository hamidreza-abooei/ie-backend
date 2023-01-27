package model

import (
	"errors"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

type Url struct {
	gorm.Model
	Id          uint   `gorm:"unique_index:index_addr_user"`
	UserId      uint   `gorm:"unique_index:index_addr_user"`
	Address     string `gorm:"unique_index:index_addr_user"`
	Threshold   int
	FailedTimes int
	Requests    []Request `gorm:"foreignkey:url_id"` // Personally I would resolve this one using DB query in order to database normalization
}

// NewURL creates a URL instance if it's address is a valid URL address
func NewURL(userID uint, address string, threshold int) (*Url, error) {
	url := new(Url)
	url.UserId = userID
	url.Threshold = threshold
	url.FailedTimes = 0

	isValid := govalidator.IsURL(address)
	if !strings.HasPrefix("http://", address) {
		address = "http://" + address
	}
	if isValid {
		//valid URL address
		url.Address = address
		return url, nil
	}
	return nil, errors.New("not a valid URL address")
}

// ShouldTriggerAlarm checks if current url's failed times is greater than it's threshold
//
// Use this function to check alarm and trigger an alarm with other functions
func (url *Url) ShouldTriggerAlarm() bool {
	return url.FailedTimes >= url.Threshold
}
