package model

import (
	"errors"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Id       uint   `gorm:"unique_index:not_null"`
	Username string `gorm:"unique_index:not_null"`
	Password string `gorm:"not null"`
	Urls     []Url  `gorm:"foreignkey:user_id"` // Personally I would resolve this one using DB query in order to database normalization
}

// NewUser creates a user with username and Hashed password
// returns error if username or password is empty
func NewUser(username, password string) (*User, error) {
	if len(password) == 0 || len(username) == 0 {
		return nil, errors.New("username of password cannot be empty")
	}
	pass, _ := HashPassword(password)
	return &User{Username: username, Password: pass}, nil
}

// HashPassword generates a hashed string from 'pass'
// returns error if 'pass' is empty
func HashPassword(pass string) (string, error) {
	if len(pass) == 0 {
		return "", errors.New("password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hash), err
}

// ValidatePassword compares 'pass' with 'users' password
// returns true if their equivalent
func (user *User) ValidatePassword(pass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)) == nil
}
