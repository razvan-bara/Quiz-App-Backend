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
		return squiz.NewAddAttemptInternalServerError().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusInternalServerError),
			Message: swag.String("couldn't save attempt"),
		})
	}

	attemptDTO = utils.ConvertAttemptModelToAttemptDTO(attempt)
	return squiz.NewAddAttemptCreated().WithPayload(attemptDTO)
}

func (handler AttemptHandler) AddAttemptAnswer(params squiz.AddAttemptAnswerParams, principal *sdto.Principal) middleware.Responder {

	ctx := context.Background()
	answerDTO := params.AttemptAnswerDTO

	// check whether question belongs to quiz
	questionArgs := &db.GetQuestionByIdAndQuizIdParams{
		ID:     answerDTO.QuestionID,
		QuizID: params.QuizID,
	}
	if _, err := handler.storage.GetQuestionByIdAndQuizId(ctx, questionArgs); err != nil {
		return squiz.NewAddAttemptAnswerBadRequest().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusBadRequest),
			Message: swag.String("question doesn't make with the quiz"),
		})
	}

	// check whether answer belongs to question
	answerArgs := &db.GetAnswerByIdAndQuestionIdParams{
		ID:         answerDTO.AnswerID,
		QuestionID: answerDTO.QuestionID,
	}
	if _, err := handler.storage.GetAnswerByIdAndQuestionId(ctx, answerArgs); err != nil {
		return squiz.NewAddAttemptAnswerBadRequest().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusBadRequest),
			Message: swag.String("answer doesn't make with the question"),
		})
	}

	createArgs := &db.CreateAttemptAnswerParams{
		AttemptID:  params.AttemptID,
		QuestionID: answerDTO.QuestionID,
		AnswerID:   answerDTO.AnswerID,
	}
	answer, err := handler.storage.CreateAttemptAnswer(ctx, createArgs)
	if err != nil {
		return squiz.NewAddAttemptAnswerInternalServerError().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusInternalServerError),
			Message: swag.String("couldn't save attempt"),
		})
	}

	answerDTO = utils.ConvertAttemptAnswerModelToAttemptAnswerDTO(answer)
	return squiz.NewAddAttemptAnswerCreated().WithPayload(answerDTO)
}

func NewAttemptHandler(storage db.Storage) *AttemptHandler {
	return &AttemptHandler{storage: storage}
}
