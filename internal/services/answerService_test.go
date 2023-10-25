package services

import (
	"context"
	"errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	mockdb "github.com/razvan-bara/VUGO-API/db/sqlc/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func getTestAnswerService(t *testing.T) (IAnswerService, *mockdb.MockStorage, *gomock.Controller) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockStorage := mockdb.NewMockStorage(ctrl)
	answerService := NewAnswerService(mockStorage)
	return answerService, mockStorage, ctrl
}

func TestAnswerService_SaveAnswer(t *testing.T) {

	t.Run("save valid answer to existing question", func(t *testing.T) {
		answerService, storage, ctrl := getTestAnswerService(t)
		defer ctrl.Finish()

		var questionID int64 = 1
		answerDTO := generateAnswerDTO()
		args := generateCreateAnswerParams(questionID, answerDTO)
		expAnswer := generateAnswer(args, questionID)

		storage.EXPECT().
			CreateAnswer(gomock.Any(), gomock.Eq(args)).
			Times(1).
			Return(expAnswer, nil)

		gotAnswer, err := answerService.SaveAnswer(context.Background(), questionID, answerDTO)
		require.NoError(t, err)

		require.Equal(t, gotAnswer.ID, expAnswer.ID)
		require.Equal(t, gotAnswer.Title, expAnswer.Title)
		require.Equal(t, gotAnswer.Correct, expAnswer.Correct)
		require.Equal(t, gotAnswer.UUID, expAnswer.UUID)
		require.Equal(t, gotAnswer.QuestionID, expAnswer.QuestionID)
		require.True(t, gotAnswer.CreatedAt.Equal(expAnswer.CreatedAt))

		require.True(t, gotAnswer.CreatedAt.After(time.Time(answerDTO.CreatedAt)))
		require.NotEqual(t, gotAnswer.ID, answerDTO.ID)
		require.NotEqual(t, gotAnswer.UUID.String(), answerDTO.UUID.String())
		require.NotEqual(t, gotAnswer.QuestionID, answerDTO.QuizQuestionID)
	})

	t.Run("return err on invalid answer.questionID", func(t *testing.T) {
		answerService, storage, ctrl := getTestAnswerService(t)
		defer ctrl.Finish()

		var questionID int64 = 0
		answerDTO := generateAnswerDTO()
		args := generateCreateAnswerParams(questionID, answerDTO)

		expErr := errors.New("foreign_key_violation")
		storage.EXPECT().
			CreateAnswer(gomock.Any(), gomock.Eq(args)).
			Times(1).
			Return(nil, expErr)

		answer, gotErr := answerService.SaveAnswer(context.Background(), questionID, answerDTO)
		require.Error(t, gotErr)
		require.Nil(t, answer)
		require.Equal(t, gotErr.Error(), expErr.Error())

	})

}

func generateAnswer(args *db.CreateAnswerParams, questionID int64) *db.Answer {
	return &db.Answer{
		ID:         1,
		Title:      args.Title,
		Correct:    args.Correct,
		QuestionID: questionID,
		UUID:       uuid.New(),
		CreatedAt:  time.Now(),
	}
}

func generateCreateAnswerParams(questionID int64, answerDTO *sdto.AnswerDTO) *db.CreateAnswerParams {
	return &db.CreateAnswerParams{
		QuestionID: questionID,
		Title:      swag.StringValue(answerDTO.Title),
		Correct:    swag.BoolValue(answerDTO.Correct),
	}
}

func generateAnswerDTO() *sdto.AnswerDTO {
	return &sdto.AnswerDTO{
		ID:             99,
		Correct:        swag.Bool(false),
		CreatedAt:      strfmt.DateTime(time.Now()),
		QuizQuestionID: 99,
		Title:          swag.String("Some answer title"),
		UUID:           strfmt.UUID(uuid.New().String()),
	}
}
