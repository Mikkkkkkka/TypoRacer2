# Config
SOURCE_NAME_SERVER=server
SOURCE_NAME_CLI_CLIENT=cli-client
SERVER_OUTPUT=bin/typo-$(SOURCE_NAME_SERVER)
CLI_CLIENT_OUTPUT=bin/typo-$(SOURCE_NAME_CLI_CLIENT)

.PHONY: all server client clean

all: server client

server:
	@echo "Building server..."
	@mkdir -p bin
	go build -o $(SERVER_OUTPUT) ./cmd/$(SOURCE_NAME_SERVER)

client:
	@echo "Building client..."
	@mkdir -p bin
	go build -o $(CLI_CLIENT_OUTPUT) ./cmd/$(SOURCE_NAME_CLI_CLIENT)

run-server: server
	@echo "Starting server..."
	./$(SERVER_OUTPUT) --port=5431

run-e2e: server client
	@echo "Starting server..."
	@mkdir -p mklogs/
	@./$(SERVER_OUTPUT) --port=5431 >./mklogs/server.log 2>&1 &
	@printf "Cli-client command: "
	@read COMMAND; echo "---------- $(CLI_CLIENT_OUTPUT) ----------"; \
	./$(CLI_CLIENT_OUTPUT) "$$COMMAND"; \
	killall typo-$(SOURCE_NAME_SERVER)

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf mklogs/
	@go clean


# Maybe in the future

# test:
# 	go test -v ./...

# docker-build:
# 	docker build -t $(BINARY_NAME_SERVER) .

help:
	@echo "Available targets:"
	@echo "  all       - Build server and client (default)"
	@echo "  server    - Build only server"
	@echo "  client    - Build only client"
	@echo "  run-server- Build and run server"
	@echo "  clean     - Remove build artifacts"