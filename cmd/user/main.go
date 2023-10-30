package main

import (
	"database/sql"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"
	_ "github.com/lib/pq"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	"github.com/razvan-bara/VUGO-API/api/user_api"
	"github.com/razvan-bara/VUGO-API/api/user_api/suser"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	"github.com/razvan-bara/VUGO-API/internal/handlers"
	"github.com/razvan-bara/VUGO-API/internal/services"
	"log"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://db_user:db_pass@localhost:5431/users?sslmode=disable"
)

func main() {

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("couldn't connect to db")
	}

	queries := db.NewSQLStorage(conn)
	userService := services.NewUserService(queries)
	userHandler := handlers.NewUserHandler(userService)

	swaggerSpec, err := loads.Analyzed(user_api.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	swaggerAPI := suser.NewUsersMicroserviceAPI(swaggerSpec)
	server := user_api.NewServer(swaggerAPI)
	defer server.Shutdown()
	server.EnabledListeners = []string{"http"}
	server.Port = 3002

	swaggerAPI.KeyAuth = func(s string) (*sdto.Principal, error) {
		// TODO: Implement auth function
		if s[len("Bearer "):] == "abc123" {
			prin := &sdto.Principal{
				ID:    0,
				Email: "",
				Exp:   strfmt.DateTime{},
			}
			return prin, nil
		}

		return nil, errors.New(401, "incorrect api key auth")
	}

	swaggerAPI.RegisterUserHandler = suser.RegisterUserHandlerFunc(userHandler.RegisterUser)
	swaggerAPI.LoginUserHandler = suser.LoginUserHandlerFunc(userHandler.AttemptLogin)

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
