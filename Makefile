quiz_drop_schema:
	docker exec -it vugo_quiz_db psql "postgresql://db_user:db_pass@localhost:5432/quizzes?sslmode=disable" -c "DROP DATABASE quizzes;"

quiz_db_up:
	migrate -path db/migrations/quiz -database "postgresql://db_user:db_pass@localhost:5432/quizzes?sslmode=disable" up

quiz_db_down:
	migrate -path db/migrations/quiz -database "postgresql://db_user:db_pass@localhost:5432/quizzes?sslmode=disable" down -all

quiz_reset_db:
	make quiz_db_down && make quiz_db_up


test:
	go test ./... -cover -v

sqlc:
	sqlc generate