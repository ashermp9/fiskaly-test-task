# Makefile for building and running the application

.PHONY: build run docker-build docker-run

APP_NAME := signer
CONFIG_PATH := ./config/local/config.yaml

build:
	@echo "Building the application..."
	go build -o $(APP_NAME) ./cmd/$(APP_NAME)/main.go

run-executable:
	@echo "Running the executable..."
	./$(APP_NAME) -config $(CONFIG_PATH)

run:
	@echo "Running the application..."
	go run ./cmd/$(APP_NAME)/main.go -config $(CONFIG_PATH)

clear:
	@echo "Removing the application binary..."
	rm ./$(APP_NAME) | true


docker-build:
	@echo "Building the Docker image..."
	docker build -t $(APP_NAME) .

docker-run:
	@echo "Running the application in Docker..."
	docker run --rm -p 8080:8080 -e CONFIG_PATH=/app/config.yaml -v $(PWD)/config/local/config.yaml:/app/config.yaml $(APP_NAME)


To interact with the API, you can use the following `curl` commands:


create-device:
	curl -X POST http://localhost:8080/api/v0/create-device \
		-H "Content-Type: application/json" \
			-d '{ \
				"id": "test-device-1", \
				"algorithm": "RSA", \
				"label": "Test Device" \
			}'

bad-device:
	curl -X POST http://localhost:8080/api/v0/create-device \
		-H "Content-Type: application/json" \
			-d '{ \
				"id": "", \
				"algorithm": "RSA", \
				"label": "Test Device" \
			}'


sign:
	curl -X POST http://localhost:8080/api/v0/sign-transaction \
    -H "Content-Type: application/json" \
    -d '{ \
        "deviceId": "test-device-1", \
        "data": "Sample data to be signed" \
    }'
