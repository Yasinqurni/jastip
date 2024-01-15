include .env
export

LOCAL_BIN:=$(CURDIR)/bin
PATH:=$(LOCAL_BIN):$(PATH)
APP_VERSION=v1.0.0
APP_PORT=8888
CONTAINER_NAME=tongsang_be

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker-compose
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-up-integration-test: ### Run docker-compose with integration test
	docker-compose up --build --abort-on-container-exit --exit-code-from integration
.PHONY: compose-up-integration-test

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

swag-v1: ### swag init
	swag init -g internal/controller/http/v1/router.go
.PHONY: swag-v1

run: swag-v1 ### swag run
	go mod tidy && go mod download && \
	DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=0 go run ./cmd/app/main.go
.PHONY: run

docker-rm-volume: ### remove docker volume
	docker volume rm go-clean-template_pg-data
.PHONY: docker-rm-volume

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-hadolint: ### check by hadolint linter
	git ls-files --exclude='Dockerfile*' --ignored | xargs hadolint
.PHONY: linter-hadolint

linter-dotenv: ### check by dotenv linter
	dotenv-linter
.PHONY: linter-dotenv

test: ### run test
	go test -v -cover -race ./internal/...
.PHONY: test

integration-test: ### run integration-test
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

mock: ### run mockgen
	mockgen -source ./internal/usecase/interfaces.go -package usecase_test > ./internal/usecase/mocks_test.go
.PHONY: mock

migrate-create:  ### create new migration
	goose -dir ./migrations create '$(MIGRATION_FILE)' sql
.PHONY: migrate-create

migrate-up: ### migration up
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir ./migrations up
.PHONY: migrate-up

migrate-down: ### migration down
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir ./migrations down
.PHONY: migrate-down

migrate-status: ### migration status
	GOOSE_DRIVER=$(GOOSE_DRIVER) GOOSE_DBSTRING=$(GOOSE_DBSTRING) goose -dir ./migrations status
.PHONY: migrate-status

delete-container:
	docker ps -q --filter name=$(CONTAINER_NAME) | xargs -r docker stop && docker ps -aq --filter name=$(CONTAINER_NAME) | xargs -r docker rm

build-prod:
	docker build -t tongsang_be:$(APP_VERSION) . --no-cache

run-prod-docker:
	@make delete-container && docker run -d -v ./internal/log:/internal/log -v ./.env:/.env -p $(APP_PORT):8080 --name $(CONTAINER_NAME) tongsang_be:$(APP_VERSION)

reload-prod-docker:
	docker restart $(CONTAINER_NAME)

docker-seed:
	docker exec -it tongsang_be-dev go run cmd/seeder/main.go $(TABLE_NAME)

gen-proto:
	protoc --go_out=./internal/entity/ --go-grpc_out=./internal/entity/ ./internal/entity/proto/*.proto

bin-deps:
	GOBIN=$(LOCAL_BIN) go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@latest

watch:
	 CompileDaemon --build="go build -o main cmd/app/main.go" --command="./main"
.PHONY: watch