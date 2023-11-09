package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/razvan-bara/VUGO-API/api/quiz_api/squiz"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	"github.com/razvan-bara/VUGO-API/internal/services"
	"net/http"
)

type QuizHandler struct {
	quizService services.IQuizService
}

func NewQuizHandler(quizService services.IQuizService) *QuizHandler {
	return &QuizHandler{quizService: quizService}
}

func (handler *QuizHandler) ProcessNewQuiz(req squiz.AddQuizParams, principal *sdto.Principal) middleware.Responder {

	if !principal.IsAdmin {
		return squiz.NewAddQuizUnauthorized().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusUnauthorized),
			Message: swag.String("unauthorised"),
		})
	}

	quiz := req.Body

	res, err := handler.quizService.ProcessNewQuiz(quiz, req.SaveMode)
	if err != nil {
		return squiz.NewAddQuizInternalServerError().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusInternalServerError),
			Message: swag.String("something went wrong while adding a quiz"),
		})
	}

	return squiz.NewAddQuizCreated().WithPayload(res)
}

func (handler *QuizHandler) ListQuizzesHandler(params squiz.ListQuizzesParams, principal *sdto.Principal) middleware.Responder {

	quizzes, err := handler.quizService.ListQuizzes(params.Status)
	if err != nil {
		return squiz.NewListQuizzesInternalServerError().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusInternalServerError),
			Message: swag.String("couldn't fetch all quizzes"),
		})
	}

	return squiz.NewListQuizzesOK().WithPayload(quizzes)
}

func (handler *QuizHandler) GetQuiz(params squiz.GetQuizParams, principal *sdto.Principal) middleware.Responder {

	quiz, err := handler.quizService.GetCompleteQuiz(params.ID)
	if err != nil {
		return squiz.NewGetQuizNotFound().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusNotFound),
			Message: swag.String("couldn't find specified quiz"),
		})
	}

	return squiz.NewGetQuizOK().WithPayload(quiz)
}

func (handler *QuizHandler) UpdateQuiz(params squiz.UpdateQuizParams, principal *sdto.Principal) middleware.Responder {

	if !principal.IsAdmin {
		return squiz.NewUpdateQuizUnauthorized().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusUnauthorized),
			Message: swag.String("unauthorised"),
		})
	}
	quiz, err := handler.quizService.UpdateCompleteQuiz(params.ID, params.Body, params.SaveMode)
	if err != nil {
		return squiz.NewGetQuizNotFound().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusNotFound),
			Message: swag.String("couldn't find specified quiz"),
		})
	}

	return squiz.NewUpdateQuizOK().WithPayload(quiz)
}

func (handler *QuizHandler) DeleteQuiz(params squiz.DeleteQuizParams, principal *sdto.Principal) middleware.Responder {

	if !principal.IsAdmin {
		return squiz.NewDeleteQuizUnauthorized().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusUnauthorized),
			Message: swag.String("unauthorised"),
		})
	}

	err := handler.quizService.DeleteQuiz(params.ID)
	if err != nil {
		return squiz.NewDeleteQuestionNotFound().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusNotFound),
			Message: swag.String("couldn't delete quiz"),
		})
	}

	return squiz.NewDeleteQuizNoContent()
}
