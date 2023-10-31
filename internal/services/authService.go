package services

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	gapi "github.com/razvan-bara/VUGO-API/api/grpc"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	"strings"
	"time"
)

const (
	bearerPrefix = "Bearer "
)

type AuthService struct {
	us IUserService
	gapi.UnimplementedAuthServiceServer
}

func NewAuthService(us IUserService) *AuthService {
	return &AuthService{us: us}
}

func (as *AuthService) extractUserFromToken(tokenString string) (*sdto.User, error) {

	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("$ecret"), nil
	})

	if err != nil {
		return nil, err
	}

	claims := jwtToken.Claims

	expTime, err := claims.GetExpirationTime()
	if !expTime.After(time.Now()) {
		return nil, errors.New("expired token")
	}

	subject, _ := claims.GetSubject()
	user, err := as.us.FindUserByEmail(subject)
	if err != nil {
		return nil, errors.New("Subject not found")
	}

	return user, nil
}

func (as *AuthService) extractTokenFromHeader(header string) (string, error) {

	header = strings.TrimSpace(header)
	if header == "" {
		return "", errors.New("Missing Bearer token")
	}

	hasPrefix := strings.HasPrefix(header, bearerPrefix)
	if !hasPrefix {
		return "", errors.New("Missing Bearer prefix")
	}

	tokenString := header[len(bearerPrefix):]
	return tokenString, nil
}

func (as *AuthService) ValidateJWTAuthorizationHeader(ctx context.Context, header *gapi.Header) (*gapi.User, error) {
	tokenString, err := as.extractTokenFromHeader(header.GetContent())
	if err != nil {
		return nil, err
	}

	user, err := as.extractUserFromToken(tokenString)
	if err != nil {
		return nil, err
	}

	return &gapi.User{Id: uint64(user.ID), Email: user.Email.String(), UUID: user.UUID.String(), IsAdmin: user.IsAdmin}, nil
}
