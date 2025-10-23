#-----------------------------------------#
###            Preparation              ###
#-----------------------------------------#

# Include .env as environment variables if file exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Define PostgreSQL DSN
POSTGRES_DSN := postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=$(POSTGRES_SSL)

#-----------------------------------------#
###           End Preparation           ###
#-----------------------------------------#


#-----------------------------------------#
###         Linting, formatting 		###
#-----------------------------------------#

.PHONY: lint_install
lint_install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.5

.PHONY: lint
lint:
	golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 ./...

#-----------------------------------------#
### 	  End Linting, formatting 		###
#-----------------------------------------#


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
###       End Database Migrations       ###
#-----------------------------------------#


#-----------------------------------------#
###             Build, Run              ###
#-----------------------------------------#

### TODO: write build targets

#-----------------------------------------#
###          End Build, Run             ###
#-----------------------------------------#


#-----------------------------------------#
###               Test                  ###
#-----------------------------------------#

### TODO: write test targets

#-----------------------------------------#
###               End Test               ###
#-----------------------------------------#
