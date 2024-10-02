APP=cmd/main.go
SCHEMA=./migrations
DB='postgres://postgres:yourpassword@127.0.0.1:5436/music_library?sslmode=disable'

build:
	docker-compose build app

run:
	docker-compose up app

local_build:
	go build -o bin/app.out $(APP)

local_run:
	go run $(APP) local_config

# Migration commands
migrate_up:
	migrate -path $(SCHEMA) -database $(DB) up

migrate_down:
	migrate -path $(SCHEMA) -database $(DB) down

# Test commands
run_test:
	go test ./... -cover
	go test -tags=e2e

lint:
	go fmt ./...
	golangci-lint run

create_test_db:
	pgpassword=yourpassword psql -h localhost -p 5436 -U postgres -tc "CREATE DATABASE postgres_test"

insert_test_data:
	pgpassword=yourpassword psql -h localhost -p 5436 -U postgres -d music_library -f ./scripts/insert_test_data.sql
