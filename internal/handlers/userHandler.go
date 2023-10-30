package handlers

import (
	"database/sql"
	"errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/lib/pq"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	"github.com/razvan-bara/VUGO-API/api/user_api/suser"
	"github.com/razvan-bara/VUGO-API/internal/services"
	"github.com/razvan-bara/VUGO-API/internal/utils"
	"net/http"
)

type UserHandler struct {
	userService services.IUserService
}

func (handler UserHandler) RegisterUser(params suser.RegisterUserParams) middleware.Responder {
	user, err := handler.userService.AddUser(params.Body)

	psErr, ok := err.(*pq.Error)
	if ok {
		if psErr.Code.Name() == "unique_violation" {
			return suser.NewRegisterUserBadRequest().WithPayload(&sdto.Error{
				Code:    swag.Int64(http.StatusBadRequest),
				Message: swag.String("pick a different email"),
			})
		}
	}

	if err != nil {
		return suser.NewRegisterUserInternalServerError().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusInternalServerError),
			Message: swag.String("couldn't register user"),
		})
	}

	return suser.NewRegisterUserOK().WithPayload(user)
}

func (handler UserHandler) AttemptLogin(params suser.LoginUserParams) middleware.Responder {
	user, err := handler.userService.AttemptLogin(params.Body)

	if errors.Is(err, sql.ErrNoRows) {
		return suser.NewLoginUserNotFound().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusNotFound),
			Message: swag.String("user with given email not found"),
		})
	}

	if err != nil {
		return suser.NewLoginUserBadRequest().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusBadRequest),
			Message: swag.String("given credentials didnt match with what's stored"),
		})
	}

	token, err := utils.GenerateJWTToken(user)
	if err != nil {
		return suser.NewLoginUserInternalServerError().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusInternalServerError),
			Message: swag.String("error while generating jwt"),
		})
	}

	return suser.NewLoginUserOK().WithPayload(&sdto.LoginResponse{
		AccessToken: token,
		User:        user,
	})
}

func NewUserHandler(userService services.IUserService) *UserHandler {
	return &UserHandler{userService: userService}
}
