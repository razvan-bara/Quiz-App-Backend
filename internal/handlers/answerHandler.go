package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/razvan-bara/VUGO-API/api/quiz_api/squiz"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	"github.com/razvan-bara/VUGO-API/internal/services"
	"net/http"
)

type AnswerHandler struct {
	answerService services.IAnswerService
}

func NewAnswerHandler(IAnswerService services.IAnswerService) *AnswerHandler {
	return &AnswerHandler{answerService: IAnswerService}
}

func (handler AnswerHandler) DeleteAnswer(params squiz.DeleteAnswerParams, principal *sdto.Principal) middleware.Responder {

	if !principal.IsAdmin {
		return squiz.NewDeleteAnswerUnauthorized().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusUnauthorized),
			Message: swag.String("unauthorised"),
		})
	}

	err := handler.answerService.DeleteAnswer(params.ID)
	if err != nil {
		return squiz.NewDeleteQuestionNotFound().WithPayload(&sdto.Error{
			Code:    swag.Int64(http.StatusNotFound),
			Message: swag.String("answer not found"),
		})
	}

	return squiz.NewDeleteAnswerNoContent()
}
