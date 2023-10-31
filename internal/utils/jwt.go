package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	"time"
)

func GenerateJWTToken(user *sdto.User) (string, error) {
	timeAfter7Days := &jwt.NumericDate{Time: time.Now().Add(time.Hour * 24 * 7)}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"sub": user.Email.String(),
		"exp": timeAfter7Days,
	})

	tokenString, err := token.SignedString([]byte("$ecret"))
	return tokenString, err

}
