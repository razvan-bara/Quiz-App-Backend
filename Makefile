quiz_db_up:
	migrate -path db/migrations/quiz -database "postgresql://db_user:db_pass@localhost:5432/quizes?sslmode=disable" up

quiz_db_down:
	migrate -path db/migrations/quiz -database "postgresql://db_user:db_pass@localhost:5432/quizes?sslmode=disable" down

sqlc:
	sqlc generate