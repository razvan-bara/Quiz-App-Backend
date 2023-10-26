package services

import (
	"context"
	"database/sql"
	"github.com/go-openapi/swag"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	"github.com/razvan-bara/VUGO-API/internal/utils"
)

type IQuestionService interface {
	SaveQuestion(ctx context.Context, quizID int64, question *sdto.QuestionDTO) (*db.Question, error)
	UpdateQuestion(dto *sdto.QuestionDTO) (*sdto.QuestionDTO, error)
	DeleteQuestion(questionID int64) error
}

type QuestionService struct {
	storage db.Storage
}

func NewQuestionService(storage db.Storage) *QuestionService {
	return &QuestionService{storage: storage}
}

func (qs *QuestionService) SaveQuestion(ctx context.Context, quizID int64, question *sdto.QuestionDTO) (*db.Question, error) {
	questionArgs := &db.CreateQuestionParams{
		QuizID: quizID,
		Title:  swag.StringValue(question.Title),
		Body: sql.NullString{
			String: question.Body,
			Valid:  false,
		},
	}

	if question.Body != "" {
		questionArgs.Body.Valid = true
	}

	savedQuestion, err := qs.storage.CreateQuestion(ctx, questionArgs)
	if err != nil {
		return nil, err
	}

	return savedQuestion, err
}

func (qs *QuestionService) UpdateQuestion(questionDTO *sdto.QuestionDTO) (*sdto.QuestionDTO, error) {

	args := &db.UpdateQuestionParams{
		ID:    questionDTO.ID,
		Title: swag.StringValue(questionDTO.Title),
		Body: sql.NullString{
			String: questionDTO.Body,
			Valid:  false,
		},
	}

	if args.Body.String != "" {
		args.Body.Valid = true
	}

	updatedQuestion, err := qs.storage.UpdateQuestion(context.Background(), args)
	if err != nil {
		return nil, err
	}

	return utils.ConvertQuestionModelToQuestionDTO(updatedQuestion), nil
}

func (qs *QuestionService) DeleteQuestion(questionID int64) error {
	return qs.storage.DeleteQuestion(context.Background(), questionID)
}
