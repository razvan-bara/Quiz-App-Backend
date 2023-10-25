package services

import (
	"context"
	"database/sql"
	"github.com/go-openapi/swag"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	"github.com/razvan-bara/VUGO-API/internal/utils"
)

type IQuizService interface {
	FindQuizById(id int64) (*db.Quiz, error)
	ProcessNewQuiz(quiz *sdto.QuizForm) (*sdto.QuizForm, error)
	SaveQuiz(ctx context.Context, quiz *sdto.QuizDTO) (*db.Quiz, error)
}

type QuizService struct {
	storage         db.Storage
	questionService IQuestionService
	answerService   IAnswerService
}

func NewQuizServiceStorage(storage db.Storage) *QuizService {
	return &QuizService{storage: storage}
}

func NewQuizService(storage db.Storage, questionService IQuestionService, answerService IAnswerService) *QuizService {
	return &QuizService{storage: storage, questionService: questionService, answerService: answerService}
}

func (qs *QuizService) FindQuizById(id int64) (*db.Quiz, error) {
	return qs.storage.GetQuiz(context.Background(), id)
}

func (qs *QuizService) ProcessNewQuiz(quizForm *sdto.QuizForm) (*sdto.QuizForm, error) {

	ctx := context.Background()

	quiz, err := qs.SaveQuiz(ctx, &quizForm.QuizDTO)
	if err != nil {
		return nil, err
	}

	res := utils.GenerateQuizResponse(quiz, len(quizForm.Questions))
	for i, questionWithAnswersDTO := range quizForm.Questions {

		question, err := qs.questionService.SaveQuestion(ctx, quiz.ID, &questionWithAnswersDTO.QuestionDTO)
		if err != nil {
			return nil, err
		}

		res.Questions[i] = utils.AddQuestionToQuizResponse(question, len(questionWithAnswersDTO.Answers))
		for j, answerDTO := range questionWithAnswersDTO.Answers {

			answer, err := qs.answerService.SaveAnswer(ctx, question.ID, answerDTO)
			if err != nil {
				return nil, err
			}

			res.Questions[i].Answers[j] = utils.ConvertAnswerModelToAnswerDTO(answer)
		}
	}

	return res, nil
}

func (qs *QuizService) SaveQuiz(ctx context.Context, quiz *sdto.QuizDTO) (*db.Quiz, error) {

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
