quiz_drop_schema:
	docker exec -it vugo_quiz_db psql "postgresql://db_user:db_pass@localhost:5432/quizzes?sslmode=disable" -c "DROP DATABASE quizzes;"

quiz_db_up:
	migrate -path db/migrations/quiz -database "postgresql://db_user:db_pass@localhost:5432/quizzes?sslmode=disable" up

quiz_db_down:
	migrate -path db/migrations/quiz -database "postgresql://db_user:db_pass@localhost:5432/quizzes?sslmode=disable" down -all

quiz_reset_db:
	make quiz_db_down && make quiz_db_up

user_db_up:
	migrate -path db/migrations/user -database "postgresql://db_user:db_pass@localhost:5431/users?sslmode=disable" up

user_db_down:
	migrate -path db/migrations/user -database "postgresql://db_user:db_pass@localhost:5431/users?sslmode=disable" down -all

user_reset_db:
	make user_db_down && make user_db_up

gen_quiz_swagger:
	swagger generate server -f ./api/quizSwagger.yml -t ./api --exclude-main -s quizApi -m /sdto -a /squiz --skip-tag-packages -P sdto.Principal

gen_user_swagger:
	swagger generate server -f ./api/userSwagger.yml -t ./api --exclude-main -s userApi -m /sdto -a /suser --skip-tag-packages -P sdto.Principal


test:
	go test ./... -cover -v

sqlc:
	sqlc generate

mockStorage:
	mockgen -destination ./db/sqlc/mock/storage.go --build_flags=--mod=mod -package mockdb github.com/razvan-bara/VUGO-API/db/sqlc Storage

mockServices:
	mockgen -destination ./internal/services/mock/services.go --build_flags=--mod=mod -package mockService github.com/razvan-bara/VUGO-API/internal/services IQuizService,IQuestionService,IAnswerService

run_quiz:
	rm quiz && go build -o quiz ./cmd/quiz/main.go && ./quiz

run_user:
	rm user && go build -o user ./cmd/user/main.go && ./user

run_all:
	make run_user & make run_quiz

generate_auth_proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./api/grpc/auth.proto