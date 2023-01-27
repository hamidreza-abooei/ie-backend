package db

import (
	"errors"
	"time"

	"github.com/hamidreza-abooei/ie-project/model"
	"github.com/jinzhu/gorm"
)

// We define a struct and start creating attributes based on it
type Store struct {
	db *gorm.DB
}

// Constructor
func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

// GetUserByUserName retrieves user from database based on it's ID
// this method loads user's Urls and Requests lists
// returns error if user was not found
func (s *Store) GetUserByUserName(username string) (*model.User, error) {
	user := new(model.User)
	// remove pre loading in the future if necessary
	if err := s.db.Preload("Urls").Preload("Urls.Requests").First(user, model.User{Username: username}).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserById retrieves a user from database with given id
// returns error if user was not found
func (s *Store) GetUserByID(id uint) (*model.User, error) {
	usr := &model.User{}
	usr.ID = id
	if err := s.db.Model(usr).Preload("Urls").Find(usr).Error; err != nil {
		return nil, err
	}
	return usr, nil
}

// GetAllUsers retrieves all users from database
func (s *Store) GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// AddUser add's a user to the database
func (s *Store) AddUser(user *model.User) error {
	return s.db.Create(user).Error
}

// AddUrl add's a Url to the database
func (s *Store) AddUrl(Url *model.Url) error {
	return s.db.Create(Url).Error
}

func (s *Store) GetAllUrls() ([]model.Url, error) {
	var Urls []model.Url
	if err := s.db.Model(&model.Url{}).Find(&Urls).Error; err != nil {
		return nil, err
	}
	return Urls, nil
}

// GetUrlById retrieves a Url from database based on it's ID
// returns error if an Url was not fount
func (s *Store) GetUrlById(id uint) (*model.Url, error) {
	Url := new(model.Url)
	if err := s.db.Preload("Requests").First(Url, id).Error; err != nil {
		return nil, err
	}
	return Url, nil
}

// GetUrlByUser retrieves Urls for this user
// returns error if nothing was found
func (s *Store) GetUrlsByUser(userID uint) ([]model.Url, error) {
	var Urls []model.Url
	if err := s.db.Model(&model.Url{}).Where("user_id == ?", userID).Find(&Urls).Error; err != nil {
		return nil, err
	}
	return Urls, nil
}

// UpdateUrl updates a Url to it's new value
func (s *Store) UpdateUrl(Url *model.Url) error {
	return s.db.Model(Url).Update(Url).Error
}

// DeleteUrl deletes a Url with it's requests from database
// returns an error if Url was not found
func (s *Store) DeleteUrl(UrlID uint) error {
	Url := &model.Url{}
	Url.ID = UrlID
	// for hard deleting user s.db.Unscoped()
	q := s.db.Model(Url).Preload("Requests").Delete(&model.Request{}, "Url_id == ?", UrlID).Delete(Url)
	if q.Error != nil {
		return q.Error
	}
	if q.RowsAffected == 0 {
		return errors.New("no rows found to delete at delete Url")
	}
	return nil
}

// DismissAlert sets "FailedTimes" value to 0 and updates it's record in database
// https://github.com/jinzhu/gorm/issues/202#issuecomment-52582525
func (s *Store) DismissAlert(UrlID uint) error {
	Url := &model.Url{}
	Url.ID = UrlID
	return s.db.Model(Url).Update("failed_times", 0).Error
}

// FetchAlerts retrieves Urls which "failed_times" is greater than it's "threshold" for given userID
func (s *Store) FetchAlerts(userID uint) ([]model.Url, error) {
	var Urls []model.Url
	if err := s.db.Model(&model.Url{}).Where("user_id == ? and failed_times >= threshold", userID).Find(Urls).Error; err != nil {
		return nil, err
	}
	return Urls, nil
}

// IncrementFailed increments failed_times of a Url
func (s *Store) IncrementFailed(Url *model.Url) error {
	Url.FailedTimes += 1
	return s.UpdateUrl(Url)
}

// AddRequest adds a request to database
func (s *Store) AddRequest(req *model.Request) error {
	return s.db.Create(req).Error
}

// GetRequestByUrl retrieves all requests for this Url
func (s *Store) GetRequestsByUrl(UrlID uint) ([]model.Request, error) {
	var requests []model.Request
	if err := s.db.Model(&model.Request{UrlId: UrlID}).Where("Url_id == ?", UrlID).Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// GetUserRequestsInPeriod retrieves requests between 2 time intervals
func (s *Store) GetUserRequestsInPeriod(UrlID uint, from, to time.Time) (*model.Url, error) {
	Url := &model.Url{}
	Url.ID = UrlID
	if err := s.db.Model(Url).Preload("Requests", "created_at >= ? and created_at <= ?", from, to).First(Url).Error; err != nil {
		return nil, err
	}
	return Url, nil
}
