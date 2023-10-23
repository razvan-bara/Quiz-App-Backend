package services

import (
	"VUGO-API/api/sdto"
	db "VUGO-API/db/sqlc"
	"VUGO-API/internal/utils"
	"context"
	"database/sql"
	"errors"
	"github.com/go-openapi/swag"
)

type IQuizService interface {
	FindQuizById(id int64) (*db.Quiz, error)
	ProcessNewQuiz(quiz *sdto.QuizForm) (*sdto.QuizForm, error)
	SaveQuiz(ctx context.Context, quiz *sdto.QuizForm) (*db.Quiz, error)
}

type QuizService struct {
	storage         db.Storage
	questionService IQuestionService
	answerService   IAnswerService
}

func NewQuizService(storage db.Storage, questionService IQuestionService, answerService IAnswerService) *QuizService {
	return &QuizService{storage: storage, questionService: questionService, answerService: answerService}
}

func (qs *QuizService) FindQuizById(id int64) (*db.Quiz, error) {
	return qs.storage.GetQuiz(context.Background(), id)
}

func (qs *QuizService) ProcessNewQuiz(quiz *sdto.QuizForm) (*sdto.QuizForm, error) {

	ctx := context.Background()

	savedQuiz, err := qs.SaveQuiz(ctx, quiz)
	if err != nil {
		return nil, errors.New("error while saving quiz")
	}

	res := utils.GenerateQuizResponse(savedQuiz, len(quiz.Questions))
	for i, question := range quiz.Questions {

		savedQuestion, err := qs.questionService.SaveQuestion(ctx, savedQuiz.ID, question)
		if err != nil {
			return nil, errors.New("error while saving question")
		}

		res.Questions[i] = utils.AddQuestionToQuizResponse(savedQuestion, len(question.Answers))
		for j, answer := range question.Answers {

			savedAnswer, err := qs.answerService.SaveAnswer(ctx, savedQuestion.ID, answer)
			if err != nil {
				return nil, errors.New("error while saving answer")
			}

			res.Questions[i].Answers[j] = utils.ConvertAnswerModelToAnswerDTO(savedAnswer)
		}
	}

	return res, nil
}

func (qs *QuizService) SaveQuiz(ctx context.Context, quiz *sdto.QuizForm) (*db.Quiz, error) {

	quizArgs := &db.CreateQuizParams{
		Title: swag.StringValue(quiz.Title),
		Description: sql.NullString{
			String: quiz.Description,
			Valid:  false,
		},
	}

	if quiz.Description != "" {
		quizArgs.Description.Valid = true
	}

	savedQuiz, err := qs.storage.CreateQuiz(ctx, quizArgs)
	if err != nil {
		return nil, err
	}

	return savedQuiz, nil
}
