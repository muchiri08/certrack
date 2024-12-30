include .envrc

.PHONY: run
run:
	@go run ./cmd/web -dsn=${CERTRACK-DSN}

.PHONY: watch
watch:
	@npx tailwindcss -i ./ui/static/main.css -o ./ui/static/output.css --watch

.PHONY: create
create:
	@echo 'creating migration files for ${name}'
	@migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: migrate/up
migrate/up:
	@echo 'Running up migrations...'
	@migrate -path=./migrations -database=${CERTRACK-DSN} up

.PHONY: build
build:
	@npx tailwindcss -o ui/static/output.css --minify
	@go build -o ./bin/certrack ./cmd/web
