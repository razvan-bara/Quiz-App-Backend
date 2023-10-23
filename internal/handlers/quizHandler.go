package handlers

import (
	"VUGO-API/api/quiz_api/squiz"
	"VUGO-API/internal/services"
	"github.com/go-openapi/runtime/middleware"
)

type QuizHandler struct {
	quizService services.IQuizService
}

func NewQuizHandler(quizService services.IQuizService) *QuizHandler {
	return &QuizHandler{quizService: quizService}
}

func (handler *QuizHandler) ProcessNewQuiz(req squiz.AddQuizParams) middleware.Responder {
	quiz := req.Body
	// add the quiz to db to get id
	res, err := handler.quizService.ProcessNewQuiz(quiz)
	if err != nil {

	}
	// add each question to db to get id
	// add each answer to db≠
	// build response

	return squiz.NewAddQuizCreated().WithPayload(res)
}
