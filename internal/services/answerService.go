package services

import (
	"context"
	"github.com/go-openapi/swag"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
)

type IAnswerService interface {
	SaveAnswer(ctx context.Context, questionID int64, answer *sdto.AnswerDTO) (*db.Answer, error)
}

type AnswerService struct {
	storage db.Storage
}

func NewAnswerService(storage db.Storage) *AnswerService {
	return &AnswerService{storage: storage}
}

func (as *AnswerService) SaveAnswer(ctx context.Context, questionID int64, answer *sdto.AnswerDTO) (*db.Answer, error) {
	answerArgs := &db.CreateAnswerParams{
		QuestionID: questionID,
		Title:      swag.StringValue(answer.Title),
		Correct:    swag.BoolValue(answer.Correct),
	}

	savedAnswer, err := as.storage.CreateAnswer(ctx, answerArgs)
	if err != nil {
		return nil, err
	}

	return savedAnswer, err
}
