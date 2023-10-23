package handlers

import (
	"VUGO-API/api/quiz_api/squiz"
	"VUGO-API/api/sdto"
	"VUGO-API/internal/services"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
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
			Message: swag.String(err.Error()),
		})
	}

	return squiz.NewAddQuizCreated().WithPayload(res)
}
