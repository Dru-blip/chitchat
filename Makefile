
migrateup:
	goose -dir internal/db/migrations sqlite3 data.db up

sqlc:
	sqlc generate

run:
	go run cmd/main.go