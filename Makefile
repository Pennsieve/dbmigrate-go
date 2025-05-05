.PHONY: help clean test test-ci package publish docker-clean vet tidy docker-image-clean clean-ci package-dbmigrate test-ci-local

LAMBDA_BUCKET ?= "pennsieve-cc-lambda-functions-use1"
WORKING_DIR   ?= "$(shell pwd)"
SERVICE_NAME  ?= "dbmigrate-go"

.DEFAULT: help

help:
	@echo "Make Help for $(SERVICE_NAME)"
	@echo ""
	@echo "make test			- run tests"
	@echo "make test-ci			- run tests in CI environment"
	@echo "make test-ci-local	- run tests in CI locally"

local-services:
	docker compose -f docker-compose.test.yml down --remove-orphans
	docker compose -f docker-compose.test.yml -f docker-compose.local.override.yml up -d local-testing

test: local-services
	go test -v -p 1 ./...

test-ci:
	docker compose -f docker-compose.test.yml down --remove-orphans
	docker compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from test

# If you want to run the tests in Docker locally and you are running Docker Desktop instead of Engine, then use this target instead of test-ci.
# It sets an env var needed by testcontainers-go to start its containers within Docker when running in Desktop. See https://golang.testcontainers.org/system_requirements/ci/dind_patterns/
test-ci-local:
	docker compose -f docker-compose.test.yml down --remove-orphans
	TESTCONTAINERS_HOST_OVERRIDE=host.docker.internal docker compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from test

# Spin down active docker containers.
docker-clean:
	#docker compose -f docker-compose.test.yml -f docker-compose.build-postgres.yml down
	docker compose -f docker-compose.test.yml down

clean: docker-clean
		rm -rf $(WORKING_DIR)/bin

clean-ci: clean

vet:
	go vet ./...

tidy:
	go mod tidy

