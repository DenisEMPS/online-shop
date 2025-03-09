CONFIG_PATH = ./config/config.yml
SECRET_KEY = asd123
APP_NAME = App

.PHONY: migrate_up migrate_down psql_run redis_run build

build:
	go build -o ./$(APP_NAME) cmd/main.go

psql_run:
	docker run --name=psql -d -p 5436:5432 -e POSTGRES_PASSWORD="qwerty" postgres:17

redis_run:
	docker run --name=redis -d -p 6379:6379 -e REDIS_PASSWORD="qwerty" redis:latest redis-server --requirepass "qwerty"

migrate_up:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up

migrate_down:
	migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' down

app_run: build
	SECRET_KEY=${SECRET_KEY} ./${APP_NAME} --config=${CONFIG_PATH}
