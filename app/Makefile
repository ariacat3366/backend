.PHONY: build
build:
	go build -o app

clean:
	rm -f app

.PHONY: unit_tests
unit_tests:
	go test -v ./...

.PHONY: integration_tests
integration_tests: ENV_FILE := integration_tests_env
integration_tests:
	ENV_FILE="${ENV_FILE}" \
	./scripts/integration_tests.sh

.PHONY: tests
tests: unit_tests

.PHONY: coverage.o
coverage.o:
	go test -cover -coverprofile=coverage.o ./...

.PHONY: coverage
coverage: coverage.o
	go tool cover -html=coverage.o
