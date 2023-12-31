package services

import (
	"context"
	"database/sql"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	mockdb "github.com/razvan-bara/VUGO-API/db/sqlc/mock"
	mockService "github.com/razvan-bara/VUGO-API/internal/services/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func getTestQuizService(t *testing.T) (IQuizService, *mockdb.MockStorage, *gomock.Controller) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockStorage := mockdb.NewMockStorage(ctrl)

	quizService := NewQuizServiceStorage(mockStorage)

	return quizService, mockStorage, ctrl
}

func TestQuizService_FindQuizById(t *testing.T) {

	t.Run("fetch existing quiz", func(t *testing.T) {
		quizService, storage, ctrl := getTestQuizService(t)
		defer ctrl.Finish()

		expQuiz := generateQuiz(nil)

		storage.EXPECT().
			GetQuiz(gomock.Any(), gomock.Eq(expQuiz.ID)).
			Times(1).
			Return(expQuiz, nil)

		gotQuiz, err := quizService.FindQuizById(expQuiz.ID)
		require.NoError(t, err)
		require.NotNil(t, gotQuiz)
		require.Equal(t, gotQuiz.Title, expQuiz.Title)
	})

	t.Run("fetch NON existing quiz", func(t *testing.T) {
		quizService, storage, ctrl := getTestQuizService(t)
		defer ctrl.Finish()

		quiz := generateQuiz(nil)
		quiz.ID = 0

		storage.EXPECT().
			GetQuiz(gomock.Any(), gomock.Eq(quiz.ID)).
			Times(1).
			Return(&db.Quiz{}, sql.ErrNoRows)

		expQuiz, err := quizService.FindQuizById(quiz.ID)
		require.Error(t, err)
		require.Empty(t, expQuiz)
	})

}

func TestQuizService_SaveQuiz(t *testing.T) {

	t.Run("save only title, description fields from dto", func(t *testing.T) {
		quizService, storage, ctrl := getTestQuizService(t)
		defer ctrl.Finish()

		quizDTO := generateQuizDTO()
		args := generateCreateQuizParams(quizDTO)
		expQuiz := generateQuiz(args)

		storage.EXPECT().
			CreateQuiz(gomock.Any(), gomock.Eq(args)).
			Times(1).
			Return(expQuiz, nil)

		gotQuiz, err := quizService.SaveQuiz(context.Background(), quizDTO)
		require.NoError(t, err)

		require.NotEmpty(t, expQuiz)
		require.Equal(t, gotQuiz.Title, expQuiz.Title)
		require.Equal(t, gotQuiz.Description, expQuiz.Description)
		require.Equal(t, gotQuiz.Attempts, expQuiz.Attempts)
		require.Equal(t, gotQuiz.UUID, expQuiz.UUID)
		require.WithinDuration(t, gotQuiz.CreatedAt, expQuiz.CreatedAt, time.Second)
		require.Equal(t, gotQuiz.PublishedAt, expQuiz.PublishedAt)

		require.NotEqual(t, gotQuiz.Attempts.Int32, quizDTO.Attempts)
		require.NotEqual(t, gotQuiz.ID, quizDTO.ID)
		require.NotEqual(t, gotQuiz.UUID.String(), quizDTO.UUID.String())
		require.True(t, gotQuiz.CreatedAt.After(time.Time(quizDTO.CreatedAt)))
	})

	t.Run("save quiz without description", func(t *testing.T) {
		quizService, storage, ctrl := getTestQuizService(t)
		defer ctrl.Finish()

		quizDTO := generateQuizDTO()
		quizDTO.Description = ""

		args := generateCreateQuizParams(quizDTO)
		args.Description = sql.NullString{
			String: "",
			Valid:  false,
		}

		expQuiz := generateQuiz(args)
		expQuiz.Description = sql.NullString{
			String: "",
			Valid:  false,
		}

		storage.EXPECT().
			CreateQuiz(gomock.Any(), gomock.Eq(args)).
			Times(1).
			Return(expQuiz, nil)

		gotQuiz, err := quizService.SaveQuiz(context.Background(), quizDTO)
		require.NoError(t, err)

		require.NotEmpty(t, gotQuiz)
		require.Empty(t, gotQuiz.Description)
	})
}

func TestQuizService_ProcessNewQuiz(t *testing.T) {

	ctrl := gomock.NewController(t)
	mockStorage := mockdb.NewMockStorage(ctrl)

	questionService := mockService.NewMockIQuestionService(ctrl)
	answerService := mockService.NewMockIAnswerService(ctrl)
	quizService := NewQuizService(mockStorage, questionService, answerService)
	defer ctrl.Finish()

	quizForm := &sdto.QuizForm{
		QuizDTO: *generateQuizDTO(),
		Questions: []*sdto.QuizFormQuestionsItems0{
			{
				QuestionDTO: *generateQuestionDTO(),
				Answers: []*sdto.AnswerDTO{
					generateAnswerDTO(),
					generateAnswerDTO(),
				},
			},
			{
				QuestionDTO: *generateQuestionDTO(),
				Answers: []*sdto.AnswerDTO{
					generateAnswerDTO(),
					generateAnswerDTO(),
					generateAnswerDTO(),
					generateAnswerDTO(),
				},
			},
		},
	}

	createQuizArgs := generateCreateQuizParams(&quizForm.QuizDTO)
	quiz := generateQuiz(createQuizArgs)
	mockStorage.
		EXPECT().
		CreateQuiz(gomock.Any(), gomock.Eq(createQuizArgs)).
		Times(1).
		Return(quiz, nil)

	for _, questionWithAnswersDTO := range quizForm.Questions {

		createQuestionArgs := generateQuestionCreateArgs(quiz.ID, &questionWithAnswersDTO.QuestionDTO)
		question := generateQuestion(createQuestionArgs)

		questionService.
			EXPECT().
			SaveQuestion(gomock.Any(), gomock.Eq(quiz.ID), gomock.Eq(&questionWithAnswersDTO.QuestionDTO)).
			Times(1).
			Return(question, nil)

		for _, answerDTO := range questionWithAnswersDTO.Answers {

			createAnswerArgs := generateCreateAnswerParams(question.ID, answerDTO)
			answer := generateAnswer(createAnswerArgs, question.ID)

			answerService.
				EXPECT().
				SaveAnswer(gomock.Any(), gomock.Eq(question.ID), gomock.Eq(answerDTO)).
				Times(1).
				Return(answer, nil)

		}
	}

	gotQuizForm, err := quizService.ProcessNewQuiz(quizForm)
	require.NoError(t, err)
	require.Equal(t, gotQuizForm.Title, quizForm.Title)
	require.Equal(t, gotQuizForm.Description, quizForm.Description)
	require.Equal(t, len(gotQuizForm.Questions), len(quizForm.Questions))
	for i, question := range gotQuizForm.Questions {
		require.NotZero(t, question.ID)
		require.NotZero(t, question.UUID)
		require.NotZero(t, question.CreatedAt)
		require.Equal(t, question.QuizID, gotQuizForm.QuizDTO.ID)

		expQuestion := quizForm.Questions[i]
		require.Equal(t, question.Body, expQuestion.Body)
		require.Equal(t, question.Title, expQuestion.Title)
		require.Equal(t, len(question.Answers), len(expQuestion.Answers))

		for j, answer := range question.Answers {

			require.NotZero(t, answer.ID)
			require.NotZero(t, answer.UUID)
			require.NotZero(t, answer.CreatedAt)
			require.Equal(t, answer.QuizQuestionID, question.ID)
			require.Equal(t, answer.Title, expQuestion.Answers[j].Title)
			require.Equal(t, answer.Correct, expQuestion.Answers[j].Correct)
		}
	}
}

func generateQuiz(args *db.CreateQuizParams) *db.Quiz {
	var quiz *db.Quiz
	if args != nil {
		quiz = &db.Quiz{
			ID:          1,
			Title:       args.Title,
			Description: args.Description,
			Attempts: sql.NullInt32{
				Int32: 0,
				Valid: true,
			},
			UUID:      uuid.New(),
			CreatedAt: time.Now(),
			PublishedAt: sql.NullTime{
				Valid: false,
			},
		}
	} else {
		quiz = &db.Quiz{
			ID:    1,
			Title: "Some random quiz title",
		}
	}
	return quiz
}

func generateQuizDTO() *sdto.QuizDTO {
	return &sdto.QuizDTO{
		ID:          2,
		Attempts:    9,
		CreatedAt:   strfmt.DateTime(time.Now()),
		Description: "Some description for dto",
		Title:       swag.String("Title for dto"),
		UUID:        strfmt.UUID(uuid.New().String()),
	}
}

func generateCreateQuizParams(quizDTO *sdto.QuizDTO) *db.CreateQuizParams {
	return &db.CreateQuizParams{
		Title: swag.StringValue(quizDTO.Title),
		Description: sql.NullString{
			String: quizDTO.Description,
			Valid:  true,
		},
	}
}
