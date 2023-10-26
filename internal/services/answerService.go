package services

import (
	"context"
	"github.com/go-openapi/swag"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	"github.com/razvan-bara/VUGO-API/internal/utils"
)

type IAnswerService interface {
	SaveAnswer(ctx context.Context, questionID int64, answer *sdto.AnswerDTO) (*db.Answer, error)
	UpdateAnswer(answer *sdto.AnswerDTO) (*sdto.AnswerDTO, error)
	DeleteAnswer(answerID int64) error
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

func (as *AnswerService) UpdateAnswer(answer *sdto.AnswerDTO) (*sdto.AnswerDTO, error) {

	args := &db.UpdateAnswerParams{
		ID:      answer.ID,
		Title:   swag.StringValue(answer.Title),
		Correct: swag.BoolValue(answer.Correct),
	}

	updateAnswer, err := as.storage.UpdateAnswer(context.Background(), args)
	if err != nil {
		return nil, err
	}

	return utils.ConvertAnswerModelToAnswerDTO(updateAnswer), nil
}

func (as *AnswerService) DeleteAnswer(answerID int64) error {
	return as.storage.DeleteAnswer(context.Background(), answerID)
}
