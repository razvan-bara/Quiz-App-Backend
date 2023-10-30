package services

import (
	"context"
	"github.com/go-openapi/swag"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	"github.com/razvan-bara/VUGO-API/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	AddUser(body *sdto.RegisterRequest) (*sdto.User, error)
}

type UserService struct {
	storage db.Storage
}

func (us UserService) AddUser(registerBody *sdto.RegisterRequest) (*sdto.User, error) {

	pass := swag.StringValue(registerBody.Password)
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	if err != nil {
		return nil, err
	}

	args := &db.CreateUserParams{
		Email:     registerBody.Email.String(),
		Password:  string(hashedPass),
		FirstName: swag.StringValue(registerBody.FirstName),
		LastName:  swag.StringValue(registerBody.LastName),
	}

	user, err := us.storage.CreateUser(context.Background(), args)
	if err != nil {
		return nil, err
	}

	return utils.ConvertUserModelToUserDTO(user), nil
}

func NewUserService(storage db.Storage) *UserService {
	return &UserService{storage: storage}
}
