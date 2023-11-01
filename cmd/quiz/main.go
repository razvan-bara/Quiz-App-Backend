package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	gapi "github.com/razvan-bara/VUGO-API/api/grpc"
	"github.com/razvan-bara/VUGO-API/api/quiz_api"
	"github.com/razvan-bara/VUGO-API/api/quiz_api/squiz"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	"github.com/razvan-bara/VUGO-API/internal/handlers"
	"github.com/razvan-bara/VUGO-API/internal/services"
	"google.golang.org/grpc"
	"log"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://db_user:db_pass@localhost:5432/quizzes?sslmode=disable"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("couldn't connect to db")
	}
	queries := db.NewSQLStorage(conn)

	questionService := services.NewQuestionService(queries)
	answerService := services.NewAnswerService(queries)
	quizService := services.NewQuizService(queries, questionService, answerService)
	quizHandler := handlers.NewQuizHandler(quizService)
	questionHandler := handlers.NewQuestionHandler(questionService)
	answerHandler := handlers.NewAnswerHandler(answerService)
	attemptHandler := handlers.NewAttemptHandler(queries)

	swaggerSpec, err := loads.Analyzed(quiz_api.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	swaggerAPI := squiz.NewQuizMicroserviceAPI(swaggerSpec)
	server := quiz_api.NewServer(swaggerAPI)
	defer server.Shutdown()
	server.EnabledListeners = []string{"http"}
	server.Port = 3000

	swaggerAPI.KeyAuth = func(s string) (*sdto.Principal, error) {

		opts := []grpc.DialOption{grpc.WithInsecure()}

		gconn, err := grpc.Dial(":4000", opts...)
		if err != nil {
			return nil, errors.New("couldn't make conn to validate token")
		}
		defer gconn.Close()

		client := gapi.NewAuthServiceClient(gconn)
		user, err := client.ValidateJWTAuthorizationHeader(
			context.Background(),
			&gapi.Header{Content: s},
		)

		if err != nil {
			return nil, err
		}

		principal := &sdto.Principal{
			ID:      int64(user.Id),
			Email:   strfmt.Email(user.Email),
			IsAdmin: user.IsAdmin,
		}

		return principal, nil
	}

	swaggerAPI.AddQuizHandler = squiz.AddQuizHandlerFunc(quizHandler.ProcessNewQuiz)
	swaggerAPI.ListQuizzesHandler = squiz.ListQuizzesHandlerFunc(quizHandler.ListQuizzesHandler)
	swaggerAPI.GetQuizHandler = squiz.GetQuizHandlerFunc(quizHandler.GetQuiz)
	swaggerAPI.UpdateQuizHandler = squiz.UpdateQuizHandlerFunc(quizHandler.UpdateQuiz)
	swaggerAPI.DeleteQuizHandler = squiz.DeleteQuizHandlerFunc(quizHandler.DeleteQuiz)
	swaggerAPI.DeleteQuestionHandler = squiz.DeleteQuestionHandlerFunc(questionHandler.DeleteQuestion)
	swaggerAPI.DeleteAnswerHandler = squiz.DeleteAnswerHandlerFunc(answerHandler.DeleteAnswer)

	swaggerAPI.AddAttemptHandler = squiz.AddAttemptHandlerFunc(attemptHandler.AddAttempt)
	swaggerAPI.AddAttemptAnswerHandler = squiz.AddAttemptAnswerHandlerFunc(attemptHandler.AddAttemptAnswer)
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
