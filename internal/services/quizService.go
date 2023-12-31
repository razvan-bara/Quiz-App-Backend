package services

import (
	"context"
	"database/sql"
	"github.com/go-openapi/swag"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	"github.com/razvan-bara/VUGO-API/internal/utils"
	"time"
)

type IQuizService interface {
	ListQuizzes(status string, page *int32, search *string) ([]*sdto.QuizDTO, error)
	FindQuizById(id int64) (*db.Quiz, error)
	GetCompleteQuiz(id int64) (*sdto.QuizForm, error)
	ProcessNewQuiz(quiz *sdto.QuizForm, saveMode string) (*sdto.QuizForm, error)
	SaveQuiz(ctx context.Context, quiz *sdto.QuizDTO, saveMode string) (*db.Quiz, error)
	UpdateCompleteQuiz(quizID int64, quizForm *sdto.QuizForm, saveMode string) (*sdto.QuizForm, error)
	UpdateQuiz(quizID int64, quiz *sdto.QuizDTO, saveMode string) (*sdto.QuizDTO, error)
	DeleteQuiz(quizID int64) error
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

func (qs *QuizService) ListQuizzes(status string, page *int32, search *string) ([]*sdto.QuizDTO, error) {
	var quizzes []*db.Quiz
	var err error

	if page != nil {
		args := &db.ListPublishedQuizzesByTitlePaginateParams{
			Lower:   "",
			Column2: swag.Int32Value(page),
		}
		if search != nil {
			args.Lower = swag.StringValue(search)
		}
		quizzes, err = qs.storage.ListPublishedQuizzesByTitlePaginate(context.Background(), args)
	} else {
		switch status {
		case "all":
			quizzes, err = qs.storage.ListQuizzes(context.Background())
		case "draft":
			quizzes, err = qs.storage.ListDraftQuizzes(context.Background())
		case "published":
			quizzes, err = qs.storage.ListPublishedQuizzes(context.Background())
		}
	}

	if err != nil {
		return nil, err
	}

	dtos := utils.ConvertQuizModelsToQuizDTOs(quizzes)
	return dtos, nil
}

func (qs *QuizService) GetCompleteQuiz(id int64) (*sdto.QuizForm, error) {
	qf := &sdto.QuizForm{
		QuizDTO:   sdto.QuizDTO{},
		Questions: nil,
	}
	ctx := context.Background()

	quiz, err := qs.storage.GetQuiz(ctx, id)
	qf.QuizDTO = *utils.ConvertQuizModelToQuizDTO(quiz)

	if err != nil {
		return nil, err
	}

	questions, err := qs.storage.ListQuestions(ctx, quiz.ID)
	if err != nil {
		return nil, err
	}

	qf.Questions = make([]*sdto.QuizFormQuestionsItems0, len(questions))
	for i, question := range questions {

		answers, err := qs.storage.ListAnswersForQuestion(ctx, question.ID)
		if err != nil {
			return nil, err
		}

		qf.Questions[i] = &sdto.QuizFormQuestionsItems0{
			QuestionDTO: *utils.ConvertQuestionModelToQuestionDTO(question),
			Answers:     make([]*sdto.AnswerDTO, len(answers)),
		}

		for j, answer := range answers {
			qf.Questions[i].Answers[j] = utils.ConvertAnswerModelToAnswerDTO(answer)
		}
	}

	return qf, nil
}

func (qs *QuizService) FindQuizById(id int64) (*db.Quiz, error) {
	return qs.storage.GetQuiz(context.Background(), id)
}

func (qs *QuizService) ProcessNewQuiz(quizForm *sdto.QuizForm, saveMode string) (*sdto.QuizForm, error) {

	ctx := context.Background()
	quiz, err := qs.SaveQuiz(ctx, &quizForm.QuizDTO, saveMode)
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

func (qs *QuizService) UpdateCompleteQuiz(quizID int64, quizForm *sdto.QuizForm, saveMode string) (*sdto.QuizForm, error) {

	quiz, err := qs.UpdateQuiz(quizID, &quizForm.QuizDTO, saveMode)
	if err != nil {
		return nil, err
	}

	quizForm.QuizDTO = *quiz
	for i, question := range quizForm.Questions {
		var questionDTO *sdto.QuestionDTO
		if question.ID == 0 {
			newQuestion, err := qs.questionService.SaveQuestion(context.Background(), quiz.ID, &question.QuestionDTO)
			if err != nil {
				return nil, err
			}
			questionDTO = utils.ConvertQuestionModelToQuestionDTO(newQuestion)
		} else {
			questionDTO, err = qs.questionService.UpdateQuestion(&question.QuestionDTO)
			if err != nil {
				return nil, err
			}
		}

		quizForm.Questions[i].QuestionDTO = *questionDTO
		for j, answer := range question.Answers {
			var answerDTO *sdto.AnswerDTO
			if answer.ID == 0 {
				newAnswer, err := qs.answerService.SaveAnswer(context.Background(), questionDTO.ID, answer)
				if err != nil {
					return nil, err
				}
				answerDTO = utils.ConvertAnswerModelToAnswerDTO(newAnswer)
			} else {
				answerDTO, err = qs.answerService.UpdateAnswer(answer)
			}

			quizForm.Questions[i].Answers[j] = answerDTO
		}
	}

	return quizForm, nil
}

func (qs *QuizService) SaveQuiz(ctx context.Context, quiz *sdto.QuizDTO, saveMode string) (*db.Quiz, error) {

	quizArgs := &db.CreateQuizParams{
		Title: swag.StringValue(quiz.Title),
		Description: sql.NullString{
			String: quiz.Description,
			Valid:  false,
		},
	}

	if saveMode == "publish" && (quiz.PublishedAt.IsZero() || quiz.PublishedAt.IsUnixZero()) {
		quizArgs.PublishedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
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

func (qs *QuizService) UpdateQuiz(quizID int64, quiz *sdto.QuizDTO, saveMode string) (*sdto.QuizDTO, error) {

	args := &db.UpdateQuizParams{
		ID:    quizID,
		Title: swag.StringValue(quiz.Title),
		Description: sql.NullString{
			String: quiz.Description,
			Valid:  false,
		},
	}

	if saveMode == "publish" && quiz.PublishedAt.IsZero() {
		args.PublishedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	if saveMode == "draft" && !quiz.PublishedAt.IsZero() {
		args.PublishedAt = sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		}
	}

	if args.Description.String != "" {
		args.Description.Valid = true
	}

	updatedQuiz, err := qs.storage.UpdateQuiz(context.Background(), args)
	if err != nil {
		return nil, err
	}

	return utils.ConvertQuizModelToQuizDTO(updatedQuiz), nil
}

func (qs *QuizService) DeleteQuiz(quizID int64) error {
	return qs.storage.DeleteQuiz(context.Background(), quizID)
}
