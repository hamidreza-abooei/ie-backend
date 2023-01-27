package model

import (
	"net/http"

	"github.com/jinzhu/gorm"
)

type Request struct {
	gorm.Model
	Id     uint `gorm:"unique_index"` // for preventing url duplication for a single user
	UrlId  uint
	Result int
}

// SendRequest sends a HTTP GET request to the url
// returns a *Request with result status code
func (url *Url) SendRequest() (*Request, error) {
	resp, err := http.Get(url.Address)
	req := new(Request)
	req.UrlId = url.ID
	if err != nil {
		return req, err
	}
	req.Result = resp.StatusCode
	return req, nil
}
