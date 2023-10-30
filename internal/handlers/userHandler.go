package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/lib/pq"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	"github.com/razvan-bara/VUGO-API/api/user_api/suser"
	"github.com/razvan-bara/VUGO-API/internal/services"
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

func NewUserHandler(userService services.IUserService) *UserHandler {
	return &UserHandler{userService: userService}
}
