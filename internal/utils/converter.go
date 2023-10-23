package utils

import (
	"VUGO-API/api/sdto"
	db "VUGO-API/db/sqlc"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func GenerateQuizResponse(quiz *db.Quiz, numOfQuestions int) *sdto.QuizForm {
	quizDTO := ConvertQuizModelToQuizDTO(quiz)
	return &sdto.QuizForm{
		QuizDTO:   *quizDTO,
		Questions: make([]*sdto.QuizFormQuestionsItems0, numOfQuestions),
	}
}

func ConvertQuizModelToQuizDTO(quiz *db.Quiz) *sdto.QuizDTO {
	return &sdto.QuizDTO{
		ID:          quiz.ID,
		CreatedAt:   strfmt.DateTime(quiz.CreatedAt),
		Description: quiz.Description.String,
		Title:       swag.String(quiz.Title),
		UUID:        strfmt.UUID(quiz.UUID.String()),
	}
}

func AddQuestionToQuizResponse(question *db.Question, numOfAnswers int) *sdto.QuizFormQuestionsItems0 {
	questionDTO := ConvertQuestionModelToQuestionDTO(question)
	return &sdto.QuizFormQuestionsItems0{
		QuestionDTO: *questionDTO,
		Answers:     make([]*sdto.AnswerDTO, numOfAnswers),
	}
}

func ConvertQuestionModelToQuestionDTO(question *db.Question) *sdto.QuestionDTO {
	return &sdto.QuestionDTO{
		ID:        question.ID,
		Body:      question.Body.String,
		CreatedAt: strfmt.DateTime(question.CreatedAt),
		QuizID:    question.QuizID,
		Title:     swag.String(question.Title),
		UUID:      strfmt.UUID(question.UUID.String()),
	}
}

func ConvertAnswerModelToAnswerDTO(answer *db.Answer) *sdto.AnswerDTO {
	return &sdto.AnswerDTO{
		ID:             answer.ID,
		Correct:        swag.Bool(answer.Correct),
		CreatedAt:      strfmt.DateTime(answer.CreatedAt),
		QuizQuestionID: answer.QuestionID,
		Title:          swag.String(answer.Title),
		UUID:           strfmt.UUID(answer.UUID.String()),
	}
}
