package handlers

import (
	"context"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/razvan-bara/VUGO-API/api/quiz_api/squiz"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	"github.com/razvan-bara/VUGO-API/internal/utils"
	"net/http"
)

type AttemptHandler struct {
	storage db.Storage
}

func (handler AttemptHandler) AddAttempt(params squiz.AddAttemptParams, principal *sdto.Principal) middleware.Responder {

	attemptDTO := params.AttemptDTO
	args := &db.CreateAttemptParams{
		QuizID: params.QuizID,
		UserID: principal.ID,
	}

	attempt, err := handler.storage.CreateAttempt(context.Background(), args)
	if err != nil {
		return squiz.NewAddQuizInternalServerError().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusInternalServerError),
			Message: swag.String("couldn't save attempt"),
		})
	}

	attemptDTO = utils.ConvertAttemptModelToAttemptDTO(attempt)
	return squiz.NewAddAttemptCreated().WithPayload(attemptDTO)
}

func NewAttemptHandler(storage db.Storage) *AttemptHandler {
	return &AttemptHandler{storage: storage}
}
