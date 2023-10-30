package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/razvan-bara/VUGO-API/api/sdto"
)

func GenerateJWTToken(user *sdto.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email.String(),
	})

	tokenString, err := token.SignedString([]byte("$ecret"))
	return tokenString, err

}
