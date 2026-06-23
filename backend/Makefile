.PHONY: server ,migrateup
server:
	go run cmd/api/main.go

migrateup:
	migrate -path migrations \
	-database "postgres://root:secret@localhost:5432/cicd?sslmode=disable" \
	up

migratedown:
	migrate -path migrations \
	-database "postgres://root:secret@localhost:5432/cicd?sslmode=disable" \
	down