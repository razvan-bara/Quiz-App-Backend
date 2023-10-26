package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/razvan-bara/VUGO-API/api/quiz_api/squiz"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	"github.com/razvan-bara/VUGO-API/internal/services"
	"net/http"
)

type QuestionHandler struct {
	questionHandler services.IQuestionService
}

func NewQuestionHandler(questionHandler services.IQuestionService) *QuestionHandler {
	return &QuestionHandler{questionHandler: questionHandler}
}

func (handler QuestionHandler) DeleteQuestion(params squiz.DeleteQuestionParams) middleware.Responder {
	err := handler.questionHandler.DeleteQuestion(params.ID)
	if err != nil {
		return squiz.NewDeleteQuestionNotFound().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusNotFound),
			Message: swag.String("quiz not found"),
		})
	}

	return squiz.NewDeleteQuestionNoContent()
}
