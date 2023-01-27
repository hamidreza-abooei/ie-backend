package common

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// Define a token sign
// I might wanna change it
var JWTSecret = []byte("SuperSecret")

// Generate a JSON Web Token with previously defined sign
func GenerateJWT(id uint) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString(JWTSecret)
	if err != nil {
		return "", err
	}
	return t, nil
}
