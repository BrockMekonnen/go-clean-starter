# Database
DB_USER ?= birukmk
DB_PASS ?= 112544
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= go-clean

# Exporting bin folder to the path for makefile
export PATH := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s)

# ~~~ Development Environment ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

up: dev-env dev-air             ## Startup / Spinup Docker Compose and air
down: docker-stop               ## Stop Docker
destroy: docker-teardown clean  ## Teardown (removes volumes, tmp files, etc...)


dev-env: ## Bootstrap Environment (with a Docker-Compose help).
	@ docker-compose up -d --build postgres

dev-env-test: dev-env ## Run application (within a Docker-Compose help)
	@ $(MAKE) image-build
	docker-compose up web

dev-air: $(AIR) ## Starts AIR (Continuous Development app).
	air

docker-stop:
	@ docker-compose down

docker-teardown:
	@ docker-compose down --remove-orphans -v

# ~~~ Code Actions ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

lint: $(GOLANGCI) ## Runs golangci-lint with predefined configuration
	@echo "Applying linter"
	golangci-lint version
	golangci-lint run -c .golangci.yaml ./...

build: ## Builds binary
	@ printf "Building application... "
	@ go build \
		-trimpath  \
		-o engine \
		./app/
	@ echo "done"

build-race: ## Builds binary (with -race flag)
	@ printf "Building application with race flag... "
	@ go build \
		-trimpath  \
		-race      \
		-o engine \
		./app/
	@ echo "done"

go-generate: $(MOCKERY) ## Runs go generate ./...
	go generate ./...

TESTS_ARGS := --format testname --jsonfile gotestsum.json.out
TESTS_ARGS += --max-fails 2
TESTS_ARGS += -- ./...
TESTS_ARGS += -test.parallel 2
TESTS_ARGS += -test.count    1
TESTS_ARGS += -test.failfast
TESTS_ARGS += -test.coverprofile coverage.out
TESTS_ARGS += -test.timeout 5s
TESTS_ARGS += -race

run-tests: $(GOTESTSUM)
	@ gotestsum $(TESTS_ARGS) -short

tests: run-tests $(TPARSE) ## Run Tests & parse details
	@cat gotestsum.json.out | $(TPARSE) -all -top -notests

# ~~~ Docker Build ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.ONESHELL:
image-build:
	@ echo "Docker Build"
	@ DOCKER_BUILDKIT=0 docker build \
		--file Dockerfile \
		--tag go-clean-starter \
		.



# ~~~ Database Migrations ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

POSTGRES_DSN := "postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)"


migrate-up: $(MIGRATE) ## Apply all (or N up) migrations.
	@ read -p "How many migrations you want to perform (default value: [all]): " N; \
	migrate -database $(POSTGRES_DSN) -path=misc/migrations up $${N:-all}

migrate-down: $(MIGRATE) ## Apply all (or N down) migrations.
	@ read -p "How many migrations you want to perform (default value: [all]): " N; \
	migrate -database $(POSTGRES_DSN) -path=misc/migrations down $${N:-all}

migrate-drop: $(MIGRATE) ## Drop everything inside the database.
	migrate -database $(POSTGRES_DSN) -path=misc/migrations drop

migrate-create: $(MIGRATE) ## Create a set of up/down migrations with a specified name.
	@ read -p "Please provide name for the migration: " Name; \
	migrate create -ext sql -dir misc/migrations $${Name}

migrate-status: $(MIGRATE) ## Show migration status.
	migrate -database $(POSTGRES_DSN) -path=misc/migrations version

# ~~~ Cleans ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

clean: clean-artifacts clean-docker

clean-artifacts: ## Removes Artifacts (*.out)
	@printf "Cleaning artifacts... "
	@rm -f *.out
	@echo "done."

clean-docker: ## Removes dangling docker images
	@ docker image prune -f