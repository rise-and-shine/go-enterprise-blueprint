# Load environment variables from .env if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Define PostgreSQL DSN
POSTGRES_DSN := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)


#-----------------------------------------#
###         Linting, formatting 		###
#-----------------------------------------#

.PHONY: lint-install
lint-install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.5

.PHONY: lint
lint:
	golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 ./...

.PHONY: fmt
fmt:
	golangci-lint fmt ./...


#-----------------------------------------#
###         Database Migrations         ###
#-----------------------------------------#

.PHONY: migrate-create
migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir "./migrations" create $$name sql

.PHONY: migrate-up
migrate-up:
	goose -dir "./migrations" -table "_migrations" postgres "$(POSTGRES_DSN)" up

.PHONY: migrate-down
migrate-down:
	goose -dir "./migrations" -table "_migrations" postgres "$(POSTGRES_DSN)" down


#-----------------------------------------#
###               Test                  ###
#-----------------------------------------#

### TODO: write test targets


#-----------------------------------------#
###             Build, Run              ###
#-----------------------------------------#

.PHONY: infra-up
infra-up:
	docker-compose -f dev-infra.yaml --profile full up -d --build

.PHONY: infra-down
infra-down:
	docker-compose -f dev-infra.yaml --profile full down --remove-orphans

.PHONY: run
run:
	go run ./cmd run
