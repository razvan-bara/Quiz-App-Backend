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

func (handler *QuizHandler) ProcessNewQuiz(req squiz.AddQuizParams) middleware.Responder {
	quiz := req.Body

	res, err := handler.quizService.ProcessNewQuiz(quiz)
	if err != nil {
		squiz.NewAddQuizInternalServerError().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusInternalServerError),
			Message: swag.String("something went wrong while adding a quiz"),
		})
	}

	return squiz.NewAddQuizCreated().WithPayload(res)
}
