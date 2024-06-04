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
	docker volume rm ${PROJECT_DIR}_pg_data
	docker volume rm ${PROJECT_DIR}_prom_data
	docker volume rm ${PROJECT_DIR}_jaeger_data
	docker image rm chat

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

.PHONY: migrate_new
migrate_new:
	cd migrations && goose create $(name) sql

.PHONY: proto
proto:
	@for dir in $(shell find . -type f -name go.mod -exec dirname {} \;); do \
		protoc --proto_path=./proto --go_out=$$dir --go-grpc_out=$$dir proto/api/auth.proto; \
		protoc --proto_path=./proto --grpc-gateway_out=$$dir \
                --grpc-gateway_opt=generate_unbound_methods=true \
                proto/api/auth.proto; \
	done
