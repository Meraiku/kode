include .env

.ONESHELL:

build:
	@go build -o ./.bin/auth_serv ./cmd/auth

run:build
	@./.bin/auth_serv

up:db
	cd ./sql/migrations;
	goose postgres $(DB_URL) up

down:
	cd ./sql/migrations;
	goose postgres $(DB_URL) down

docker:build
	docker build . -t $(I_PATH)/$(I_NAME)
	docker push $(I_PATH)/$(I_NAME)
	docker compose up

db:
	./scripts/db.sh