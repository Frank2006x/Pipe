POSTGRES_DB=postgres://root:secret@localhost:5432/cicd?sslmode=disable
.PHONY: server ,migrateup
server:
	go run cmd/api/main.go

migrate-up:
	migrate -path db/migrations -database "$(POSTGRES_DB)" up

migrate-down:
	migrate -path db/migrations -database "$(POSTGRES_DB)" down 1

migrate-version:
	migrate -path db/migrations -database "$(POSTGRES_DB)" version

migrate-force:
	migrate -path db/migrations -database "$(POSTGRES_DB)" force 1