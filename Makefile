# Variables
BUILD_DIR := .build
BIN_NAME := api

build:
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BIN_NAME) ./cmd/api

run:
	@go run cmd/api/main.go

migrate:
	@go run cmd/migrate/main.go

clean:
	@rm -rf $(BUILD_DIR)

.PHONY: build clean
