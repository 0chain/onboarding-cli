BINARY_NAME=keygen

## build: Build binary
build:
	@echo "Building..."
	go build -ldflags="-s -w" -o bin/${BINARY_NAME} ./
	@echo "Built!"
