quiz_drop_schema:
	docker exec -it vugo_quiz_db psql "postgresql://db_user:db_pass@localhost:5432/quizzes?sslmode=disable" -c "DROP DATABASE quizzes;"

quiz_db_up:
	migrate -path db/migrations/quiz -database "postgresql://db_user:db_pass@localhost:5432/quizzes?sslmode=disable" up

quiz_db_down:
	migrate -path db/migrations/quiz -database "postgresql://db_user:db_pass@localhost:5432/quizzes?sslmode=disable" down -all

quiz_reset_db:
	make quiz_db_down && make quiz_db_up

gen_quiz_swagger:
	swagger generate server -f ./api/quizSwagger.yml -t ./api --exclude-main -s quizApi -m /sdto -a /squiz --skip-tag-packages


test:
	go test ./... -cover -v

sqlc:
	sqlc generate

run_quiz:
	go run ./cmd/quiz/main.go