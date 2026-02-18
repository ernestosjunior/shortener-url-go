run-dev:
	go run cmd/api/main.go

sqlc-generate:
	sqlc generate -f internal/store/mysql/sqlc.yml