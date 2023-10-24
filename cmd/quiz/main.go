package main

import (
	"database/sql"
	"github.com/go-openapi/loads"
	_ "github.com/golang/mock/mockgen/model"
	_ "github.com/lib/pq"
	"github.com/razvan-bara/VUGO-API/api/quiz_api"
	"github.com/razvan-bara/VUGO-API/api/quiz_api/squiz"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	"github.com/razvan-bara/VUGO-API/internal/handlers"
	"github.com/razvan-bara/VUGO-API/internal/services"
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
	queries := db.New(conn)

	questionService := services.NewQuestionService(queries)
	answerService := services.NewAnswerService(queries)
	quizService := services.NewQuizService(queries, questionService, answerService)
	quizHandler := handlers.NewQuizHandler(quizService)

	swaggerSpec, err := loads.Analyzed(quiz_api.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	swaggerAPI := squiz.NewQuizMicroserviceAPI(swaggerSpec)
	server := quiz_api.NewServer(swaggerAPI)
	defer server.Shutdown()
	server.EnabledListeners = []string{"http"}
	server.Port = 3000

	swaggerAPI.AddQuizHandler = squiz.AddQuizHandlerFunc(quizHandler.ProcessNewQuiz)

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
