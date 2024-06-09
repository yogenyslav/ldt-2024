include .env

.PHONY: lint
lint:
	@echo "Starting linter"
	@for dir in $(shell find . -type f -name go.mod -exec dirname {} \;); do \
		echo "Running linter in $$dir"; \
		cd $$dir && golangci-lint run --config $(PROJECT_DIR)/.golangci.yml && cd ..; \
	done

.PHONY: docker_up
docker_up:
	docker compose up -d --build

.PHONY: docker_down
docker_down:
	docker compose down

.PHONY: docker_remove
docker_remove: docker_down
	docker volume rm ${BASE_IMAGE}_pg_data
	docker volume rm ${BASE_IMAGE}_prom_data
	docker volume rm ${BASE_IMAGE}_jaeger_data
	docker volume rm ${BASE_IMAGE}_redis_data
	docker volume rm ${BASE_IMAGE}_redis_conf
	docker image rm chat
	docker image rm api
	docker image rm bot
	docker image rm prompter

.PHONY: docker_restart
docker_restart: docker_down docker_up

.PHONY: docker_purge_restart
docker_purge_restart: docker_remove docker_up

.PHONY: migrate_up
migrate_up:
	cd migrations && goose postgres "user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable \
		host=${POSTGRES_HOST} port=${POSTGRES_PORT}" up

.PHONY: migrate_down
migrate_down:
	cd migrations && goose postgres "user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable \
		host=localhost port=${POSTGRES_PORT}" down

.PHONY: migrate_up_test
migrate_up_test:
	cd migrations && goose postgres "user=${POSTGRES_TEST_USER} \
		password=${POSTGRES_TEST_PASSWORD} dbname=${POSTGRES_TEST_DB} sslmode=disable \
		host=${POSTGRES_TEST_HOST} port=${POSTGRES_TEST_PORT}" up

.PHONY: migrate_down_test
migrate_down_test:
	cd migrations && goose postgres "user=${POSTGRES_TEST_USER} \
		password=${POSTGRES_TEST_PASSWORD} dbname=${POSTGRES_TEST_DB} sslmode=disable \
		host=localhost port=${POSTGRES_TEST_PORT}" down

.PHONY: migrate_new
migrate_new:
	cd migrations && goose create $(name) sql

.PHONY: proto
proto:
	@for dir in $(shell find . -type f -name go.mod -exec dirname {} \;); do \
		protoc --proto_path=./proto --go_out=$$dir --go-grpc_out=$$dir proto/api/auth.proto proto/api/prompter.proto; \
	done
	@protoc --proto_path=./proto --grpc-gateway_out=./api \
                    --grpc-gateway_opt=generate_unbound_methods=true \
                    proto/api/auth.proto proto/api/prompter.proto --openapiv2_out ./api/docs/third_party/OpenAPI
	@python -m grpc_tools.protoc -Iproto --python_out=prompter --grpc_python_out=prompter \
 					proto/api/prompter.proto

.PHONY: tests
tests:
	@for dir in $(shell find . -type f -name go.mod -exec dirname {} \;); do \
		cd $$dir && go test -v ./... -cover -tags=integration && cd ..; \
	done

.PHONY: swag
swag:
	cd ./chat && swag init -g cmd/server/main.go -o ./docs
