package utils

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/strfmt/conv"
	"github.com/go-openapi/swag"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
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
		Attempts:    quiz.Attempts.Int32,
		CreatedAt:   strfmt.DateTime(quiz.CreatedAt),
		Description: quiz.Description.String,
		PublishedAt: strfmt.DateTime(quiz.PublishedAt.Time),
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

func ConvertQuizModelsToQuizDTOs(questions []*db.Quiz) []*sdto.QuizDTO {
	q := make([]*sdto.QuizDTO, len(questions))
	for i := 0; i < len(q); i++ {
		q[i] = ConvertQuizModelToQuizDTO(questions[i])
	}
	return q
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

func ConvertAttemptsModelToAttemptsDTO(attempts []*db.Attempt) []*sdto.AttemptDTO {
	a := make([]*sdto.AttemptDTO, len(attempts))
	for i := 0; i < len(attempts); i++ {
		a[i] = ConvertAttemptModelToAttemptDTO(attempts[i])
	}

	return a
}

func ConvertAttemptModelToAttemptDTO(attempt *db.Attempt) *sdto.AttemptDTO {
	return &sdto.AttemptDTO{
		ID:        attempt.ID,
		UUID:      strfmt.UUID(attempt.UUID.String()),
		CreatedAt: strfmt.DateTime(attempt.CreatedAt),
		QuizID:    attempt.QuizID,
		Score:     int32(attempt.Score),
		Status:    int32(attempt.Status),
		UserID:    attempt.UserID,
	}
}

func ConvertUserModelToUserDTO(user *db.User) *sdto.User {
	return &sdto.User{
		ID:        user.ID,
		UUID:      strfmt.UUID(user.UUID.String()),
		CreatedAt: strfmt.DateTime(user.CreatedAt),
		Email:     conv.Email(strfmt.Email(user.Email)),
		FirstName: swag.String(user.FirstName),
		LastName:  swag.String(user.LastName),
		IsAdmin:   user.IsAdmin,
	}
}

func ConvertAttemptAnswerModelToAttemptAnswerDTO(answer *db.AttemptAnswer) *sdto.AttemptAnswerDTO {
	return &sdto.AttemptAnswerDTO{
		ID:         answer.ID,
		UUID:       strfmt.UUID(answer.UUID.String()),
		AnswerID:   answer.AnswerID,
		AttemptID:  answer.AttemptID,
		CreatedAt:  strfmt.DateTime(answer.CreatedAt),
		QuestionID: answer.QuestionID,
	}
}
