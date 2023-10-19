package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func assertQuizWithArgs(t *testing.T, arg CreateQuizParams, quiz Quiz) {
	t.Helper()

	require.Equal(t, arg.Title, quiz.Title)
	require.Equal(t, arg.Description, quiz.Description)
	require.NotZero(t, quiz.ID)
	require.NotZero(t, quiz.CreatedAt)
}

func assertQuizzesEqual(t *testing.T, quiz, target Quiz) {
	t.Helper()
	require.NotEmpty(t, quiz)
	require.NotEmpty(t, target)
	require.Equal(t, quiz.Title, target.Title)
	require.Equal(t, quiz.Description, target.Description)
	require.WithinDuration(t, quiz.CreatedAt.Time, target.CreatedAt.Time, time.Second)
	require.Equal(t, quiz.UUID, target.UUID)
	require.WithinDuration(t, quiz.PublishedAt.Time, target.PublishedAt.Time, time.Second)
	require.Equal(t, quiz.Attempts, target.Attempts)
}

func addNewQuiz(t *testing.T) (Quiz, error) {
	t.Helper()
	arg := CreateQuizParams{
		Title: "Golang Quiz",
		Description: sql.NullString{
			String: "A quiz for the Golang programming language",
			Valid:  true,
		},
	}

	quiz, err := testQueries.CreateQuiz(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, quiz)
	assertQuizWithArgs(t, arg, quiz)
	return quiz, err
}

func TestCreateQuiz(t *testing.T) {

	//t.Run("create quiz with empty title", func(t *testing.T) {
	//	arg := CreateQuizParams{}
	//
	//	quiz, err := testQueries.CreateQuiz(context.Background(), arg)
	//	require.Error(t, err)
	//	require.Empty(t, quiz)
	//})

	t.Run("create quiz", func(t *testing.T) {
		addNewQuiz(t)
	})

}

func TestGetQuiz(t *testing.T) {
	quiz, _ := addNewQuiz(t)
	target, err := testQueries.GetQuiz(context.Background(), quiz.ID)

	require.NoError(t, err)
	assertQuizzesEqual(t, quiz, target)
}

func TestUpdateQuiz(t *testing.T) {
	quiz, _ := addNewQuiz(t)

	arg := UpdateQuizParams{
		ID:    quiz.ID,
		Title: "some random title change",
		Description: sql.NullString{
			String: "some random description change",
			Valid:  true,
		},
	}

	target, err := testQueries.UpdateQuiz(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, quiz)
	require.NotEmpty(t, target)
	require.NotEqual(t, quiz.Title, target.Title)
	require.NotEqual(t, quiz.Description, target.Description)
	require.Equal(t, arg.Title, target.Title)
	require.Equal(t, arg.Description, target.Description)
}

func TestDeleteQuiz(t *testing.T) {
	quiz, _ := addNewQuiz(t)

	err := testQueries.DeleteQuiz(context.Background(), quiz.ID)
	require.NoError(t, err)

	target, err := testQueries.GetQuiz(context.Background(), quiz.ID)
	require.Error(t, err)
	require.Empty(t, target)
}
