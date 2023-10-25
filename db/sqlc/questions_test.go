package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"testing"
)

func assertQuestionWithArgs(t *testing.T, arg *CreateQuestionParams, question *Question) {
	t.Helper()

	require.NotZero(t, question.ID)
	require.Equal(t, arg.Title, question.Title)
	require.Equal(t, arg.Body.String, question.Body.String)
	require.Equal(t, arg.QuizID, question.QuizID)
	require.NotEmpty(t, question.UUID)
	require.NotZero(t, question.CreatedAt)
}

func insertQuestion(t *testing.T, quizID int64) *Question {
	arg := &CreateQuestionParams{
		QuizID: quizID,
		Title:  "Some title question",
		Body: sql.NullString{
			String: "Some question body",
			Valid:  true,
		},
	}

	question, err := testQueries.CreateQuestion(context.Background(), arg)
	require.NoError(t, err)
	assertQuestionWithArgs(t, arg, question)
	return question
}

func addNewQuestionForQuiz(t *testing.T, quizID int64) *Question {
	question := insertQuestion(t, quizID)
	return question
}

func addNewQuestion(t *testing.T) *Question {
	quiz, err := addNewQuiz(t)
	require.NoError(t, err)
	question := insertQuestion(t, quiz.ID)
	return question
}

func TestCreateQuestion(t *testing.T) {

	t.Run("valid question creation", func(t *testing.T) {
		addNewQuestion(t)
	})

	t.Run("invalid quiz foreign key for question", func(t *testing.T) {
		arg := &CreateQuestionParams{
			QuizID: 999999,
			Title:  "Some title question",
			Body: sql.NullString{
				String: "Some question body",
				Valid:  true,
			},
		}

		question, err := testQueries.CreateQuestion(context.Background(), arg)
		require.Error(t, err)
		require.Empty(t, question)

		var pqErr *pq.Error
		ok := errors.As(err, &pqErr)
		require.True(t, ok)
		require.Equal(t, pqErr.Code.Name(), "foreign_key_violation")
	})

}

func TestGetQuestion(t *testing.T) {

	t.Run("questions exists in db", func(t *testing.T) {

		question := addNewQuestion(t)
		target, err := testQueries.GetQuestion(context.Background(), question.ID)
		require.NoError(t, err)
		require.Equal(t, question, target)

	})

	t.Run("question doesnt exist in db", func(t *testing.T) {

		var questionID int64 = 99999
		target, err := testQueries.GetQuestion(context.Background(), questionID)
		require.ErrorIs(t, err, sql.ErrNoRows)
		require.Empty(t, target)

	})

}

func TestListQuestions(t *testing.T) {

	t.Run("list question for valid quiz", func(t *testing.T) {
		quiz, _ := addNewQuiz(t)

		noOfQuestion := 5
		for i := 0; i < noOfQuestion; i++ {
			addNewQuestionForQuiz(t, quiz.ID)
		}

		questions, err := testQueries.ListQuestions(context.Background(), quiz.ID)
		require.NoError(t, err)
		require.Len(t, questions, noOfQuestion)
	})

	t.Run("list nothing for quiz that doesnt exist", func(t *testing.T) {

		var quizID int64 = 99999

		questions, err := testQueries.ListQuestions(context.Background(), quizID)
		require.NoError(t, err)
		require.Len(t, questions, 0)
	})

}

func TestUpdateQuestion(t *testing.T) {

	t.Run("valid question update", func(t *testing.T) {

		ctx := context.Background()
		question := addNewQuestion(t)

		args := &UpdateQuestionParams{
			ID:    question.ID,
			Title: "Updated question title",
			Body: sql.NullString{
				String: "Updated question body",
				Valid:  true,
			},
		}

		target, err := testQueries.UpdateQuestion(ctx, args)
		require.NoError(t, err)
		require.Equal(t, target.Title, args.Title)
		require.Equal(t, target.Body, args.Body)
		require.Equal(t, target.UUID, question.UUID)
		require.Equal(t, target.QuizID, question.QuizID)
		require.Equal(t, target.CreatedAt, question.CreatedAt)

	})
}

func TestDeleteQuestion(t *testing.T) {

	t.Run("delete existing question in db", func(t *testing.T) {

		question := addNewQuestion(t)
		ctx := context.Background()

		err := testQueries.DeleteQuestion(ctx, question.ID)
		require.NoError(t, err)

		q, err := testQueries.GetQuestion(ctx, question.ID)
		require.Empty(t, q)
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("delete not existing question in db", func(t *testing.T) {
		var questionID int64 = 99999

		err := testQueries.DeleteQuestion(context.Background(), questionID)
		require.NoError(t, err)
	})
}

func TestDeleteAllQuestionsForQuiz(t *testing.T) {

	quiz, _ := addNewQuiz(t)
	ctx := context.Background()

	noOfQuestion := 2
	for i := 0; i < noOfQuestion; i++ {
		addNewQuestionForQuiz(t, quiz.ID)
	}

	err := testQueries.DeleteQuestions(ctx, quiz.ID)
	require.NoError(t, err)

	questions, err := testQueries.ListQuestions(ctx, quiz.ID)
	require.NoError(t, err)
	require.Len(t, questions, 0)

}
