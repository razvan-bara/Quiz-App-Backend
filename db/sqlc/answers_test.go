package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func assertEqualAnswerWithArgs(t *testing.T, answer *Answer, args *CreateAnswerParams) {
	require.NotZero(t, answer.ID)
	require.Equal(t, answer.Title, args.Title)
	require.Equal(t, answer.Correct, args.Correct)
	require.Equal(t, answer.QuestionID, args.QuestionID)
	require.NotEmpty(t, answer.UUID)
	require.NotEmpty(t, answer.CreatedAt)
}

func addNewAnswer(t *testing.T, question *Question) *Answer {

	if question == nil {
		question = addNewQuestion(t)
	}

	args := &CreateAnswerParams{
		QuestionID: question.ID,
		Title:      "Some questions",
		Correct:    false,
	}

	answer, err := testQueries.CreateAnswer(context.Background(), args)
	require.NoError(t, err)
	assertEqualAnswerWithArgs(t, answer, args)
	return answer
}

func TestCreateAnswer(t *testing.T) {
	question := addNewQuestion(t)

	t.Run("create valid answer", func(t *testing.T) {
		addNewAnswer(t, question)
	})

	t.Run("assert default value for answer.Correct field", func(t *testing.T) {
		args := &CreateAnswerParams{
			QuestionID: question.ID,
			Title:      "Some questions",
		}

		answer, err := testQueries.CreateAnswer(context.Background(), args)
		require.NoError(t, err)
		require.False(t, answer.Correct)
	})

	t.Run("create answer for non-existing quiz", func(t *testing.T) {
		args := &CreateAnswerParams{
			QuestionID: 0,
			Title:      "Some questions",
		}

		answer, err := testQueries.CreateAnswer(context.Background(), args)
		require.Error(t, err)
		require.Empty(t, answer)
	})

}

func TestGetAnswer(t *testing.T) {

	t.Run("fetch existing answer", func(t *testing.T) {
		answer := addNewAnswer(t, nil)

		target, err := testQueries.GetAnswer(context.Background(), answer.ID)
		require.NoError(t, err)
		require.Equal(t, answer, target)
	})

	t.Run("fetch NOT existing answer", func(t *testing.T) {
		var invalidAnswerID int64 = 0
		target, err := testQueries.GetAnswer(context.Background(), invalidAnswerID)
		require.Error(t, err)
		require.Empty(t, target)
	})
}

func TestListAnswersForQuestion(t *testing.T) {

	t.Run("get answers for existing question", func(t *testing.T) {
		question := addNewQuestion(t)
		noOfAnswers := 4
		for i := 0; i < noOfAnswers; i++ {
			addNewAnswer(t, question)
		}

		answers, err := testQueries.ListAnswersForQuestion(context.Background(), question.ID)
		require.NoError(t, err)
		require.Len(t, answers, noOfAnswers)
	})

	t.Run("get answers for NON existing question", func(t *testing.T) {
		var invalidQuestionID int64 = 0
		answers, err := testQueries.ListAnswersForQuestion(context.Background(), invalidQuestionID)
		require.NoError(t, err)
		require.Len(t, answers, 0)
	})

}

func TestUpdateAnswer(t *testing.T) {

	t.Run("update existing answer fields", func(t *testing.T) {
		answer := addNewAnswer(t, nil)

		args := &UpdateAnswerParams{
			ID:      answer.ID,
			Title:   "Some changed title",
			Correct: true,
		}

		target, err := testQueries.UpdateAnswer(context.Background(), args)
		require.NoError(t, err)
		require.Equal(t, args.Title, target.Title)
		require.Equal(t, args.Correct, target.Correct)
		require.Equal(t, answer.ID, target.ID)
	})

	t.Run("update NON existing answer", func(t *testing.T) {

		var invalidAnswerID int64 = 0
		args := &UpdateAnswerParams{
			ID:      invalidAnswerID,
			Title:   "Some changed title",
			Correct: true,
		}

		target, err := testQueries.UpdateAnswer(context.Background(), args)
		require.Error(t, err)
		require.Empty(t, target)
	})

}

func TestDeleteAnswer(t *testing.T) {

	t.Run("delete existing answer", func(t *testing.T) {
		answer := addNewAnswer(t, nil)

		ctx := context.Background()
		err := testQueries.DeleteAnswer(ctx, answer.ID)
		require.NoError(t, err)

		target, err := testQueries.GetAnswer(ctx, answer.ID)
		require.Error(t, err)
		require.Empty(t, target)
	})

	t.Run("delete NOT existing answer", func(t *testing.T) {

		ctx := context.Background()
		err := testQueries.DeleteAnswer(ctx, 0)
		require.NoError(t, err)
	})
}

func TestDeleteAnswersForQuestion(t *testing.T) {

	t.Run("delete answers for existing question", func(t *testing.T) {
		question := addNewQuestion(t)
		noOfAnswers := 4
		for i := 0; i < noOfAnswers; i++ {
			addNewAnswer(t, question)
		}

		ctx := context.Background()
		err := testQueries.DeleteAnswersForQuestion(ctx, question.ID)
		require.NoError(t, err)

		answers, err := testQueries.ListAnswersForQuestion(ctx, question.ID)
		require.NoError(t, err)
		require.Len(t, answers, 0)
	})

	t.Run("delete answers for NON existing question", func(t *testing.T) {

		var invalidQuestionID int64 = 0
		err := testQueries.DeleteAnswersForQuestion(context.Background(), invalidQuestionID)
		require.NoError(t, err)

	})
}
