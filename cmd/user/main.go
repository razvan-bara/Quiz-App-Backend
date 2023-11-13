package main

import (
	"context"
	"database/sql"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"
	_ "github.com/lib/pq"
	gapi "github.com/razvan-bara/VUGO-API/api/grpc"
	"github.com/razvan-bara/VUGO-API/api/sdto"
	"github.com/razvan-bara/VUGO-API/api/user_api"
	"github.com/razvan-bara/VUGO-API/api/user_api/suser"
	db "github.com/razvan-bara/VUGO-API/db/sqlc"
	"github.com/razvan-bara/VUGO-API/internal/handlers"
	servers "github.com/razvan-bara/VUGO-API/internal/server"
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
	authService := services.NewAuthService(userService)
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

	swaggerAPI.KeyAuth = func(tokenString string) (*sdto.Principal, error) {

		user, err := authService.ValidateJWTAuthorizationHeader(context.Background(), &gapi.Header{Content: tokenString})
		if err != nil {
			return nil, errors.New(401, "incorrect api key auth")
		}

		principal := &sdto.Principal{
			ID:      int64(user.Id),
			Email:   strfmt.Email(user.Email),
			IsAdmin: user.IsAdmin,
		}

		return principal, nil
	}
	swaggerAPI.RegisterUserHandler = suser.RegisterUserHandlerFunc(userHandler.RegisterUser)
	swaggerAPI.LoginUserHandler = suser.LoginUserHandlerFunc(userHandler.AttemptLogin)
	swaggerAPI.GetUserDetailsHandler = suser.GetUserDetailsHandlerFunc(userHandler.GetUserDetails)

	server.ConfigureAPI()
	go func() {
		if err := server.Serve(); err != nil {
			log.Println(err)
		}
	}()
	servers.LoadAuthGRPCServer(authService)

	select {}
}
