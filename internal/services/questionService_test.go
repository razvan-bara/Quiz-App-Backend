package services

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	mockdb "github.com/razvan-bara/VUGO-API/db/sqlc/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func getTestQuestionService(t *testing.T) (IQuestionService, *mockdb.MockStorage, *gomock.Controller) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockStorage := mockdb.NewMockStorage(ctrl)
	questionService := NewQuestionService(mockStorage)
	return questionService, mockStorage, ctrl
}

func TestQuestionService_SaveQuestion(t *testing.T) {

	t.Run("save valid question", func(t *testing.T) {
		questionService, storage, ctrl := getTestQuestionService(t)
		defer ctrl.Finish()

		var quizID int64 = 1

		questionDTO := generateQuestionDTO()
		args := generateQuestionCreateArgs(quizID, questionDTO)
		expQuestion := generateQuestion(args)

		storage.EXPECT().
			CreateQuestion(gomock.Any(), gomock.Eq(args)).
			Times(1).
			Return(expQuestion, nil)

		gotQuestion, err := questionService.SaveQuestion(context.Background(), quizID, questionDTO)
		require.NoError(t, err)

		require.Equal(t, gotQuestion.Title, expQuestion.Title)
		require.Equal(t, gotQuestion.Body, expQuestion.Body)
		require.Equal(t, gotQuestion.QuizID, expQuestion.QuizID)

		require.NotEqual(t, gotQuestion.ID, questionDTO.ID)
		require.NotEqual(t, gotQuestion.UUID.String(), questionDTO.UUID.String())
		require.True(t, gotQuestion.CreatedAt.After(time.Time(questionDTO.CreatedAt)))
	})

	t.Run("save question with empty description", func(t *testing.T) {
		questionService, storage, ctrl := getTestQuestionService(t)
		defer ctrl.Finish()

		var quizID int64 = 1
		questionDTO := generateQuestionDTO()
		questionDTO.Body = ""

		args := generateQuestionCreateArgs(quizID, questionDTO)

		args.Body = sql.NullString{}
		expQuestion := generateQuestion(args)

		storage.EXPECT().
			CreateQuestion(gomock.Any(), gomock.Eq(args)).
			Times(1).
			Return(expQuestion, nil)

		gotQuestion, err := questionService.SaveQuestion(context.Background(), quizID, questionDTO)
		require.NoError(t, err)

		require.Empty(t, gotQuestion.Body)
	})

	t.Run("return err on invalid question.quizID", func(t *testing.T) {
		questionService, storage, ctrl := getTestQuestionService(t)
		defer ctrl.Finish()

		var quizID int64 = 0
		questionDTO := generateQuestionDTO()
		args := generateQuestionCreateArgs(quizID, questionDTO)

		expErr := errors.New("foreign_key_violation")
		storage.EXPECT().
			CreateQuestion(gomock.Any(), gomock.Eq(args)).
			Times(1).
			Return(nil, expErr)

		gotQuestion, gotErr := questionService.SaveQuestion(context.Background(), quizID, questionDTO)
		require.Error(t, gotErr)
		require.Nil(t, gotQuestion)
		require.Equal(t, expErr.Error(), gotErr.Error())
	})

}

func generateQuestionDTO() *sdto.QuestionDTO {
	return &sdto.QuestionDTO{
		ID:        99,
		Body:      "Some question body",
		CreatedAt: strfmt.DateTime(time.Now()),
		QuizID:    99,
		Title:     swag.String("Some title for question"),
		UUID:      strfmt.UUID(uuid.New().String()),
	}
}

func generateQuestionCreateArgs(quizID int64, dto *sdto.QuestionDTO) *db.CreateQuestionParams {
	return &db.CreateQuestionParams{
		QuizID: quizID,
		Title:  swag.StringValue(dto.Title),
		Body: sql.NullString{
			String: dto.Body,
			Valid:  true,
		},
	}
}

func generateQuestion(args *db.CreateQuestionParams) *db.Question {
	return &db.Question{
		ID:        1,
		Title:     args.Title,
		Body:      args.Body,
		QuizID:    args.QuizID,
		UUID:      uuid.New(),
		CreatedAt: time.Now(),
	}
}
