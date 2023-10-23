package main

import (
	"VUGO-API/api/quiz_api"
	"VUGO-API/api/quiz_api/squiz"
	db "VUGO-API/db/sqlc"
	"VUGO-API/internal/handlers"
	"VUGO-API/internal/services"
	"database/sql"
	"github.com/go-openapi/loads"
	_ "github.com/lib/pq"
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
