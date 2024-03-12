# Makefile for building Go application for multiple platforms

GO_FILES="./cmd/example"
BINARY_NAME="de-dribbble"
BIN_DIR=".builds"
APP_DATA="./df_data"

.PHONY: go-build go-run \
	windows linux darwin clean wipe-data \
	lint test test-report

go-build: windows linux darwin

go-run:
	@go run $(GO_FILES)

windows:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(BIN_DIR)
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 \
		go build -ldflags="-s -w" \
		-o $(BIN_DIR)/$(BINARY_NAME)_windows_amd64.exe $(GO_FILES)

linux:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(BIN_DIR)
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
		go build -ldflags="-s -w" \
		-o $(BIN_DIR)/$(BINARY_NAME)_linux_amd64 $(GO_FILES)

darwin:
	@echo "Building for macOS (amd64)..."
	@mkdir -p $(BIN_DIR)
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
		go build -ldflags="-s -w" \
		-o $(BIN_DIR)/$(BINARY_NAME)_darwin_amd64 $(GO_FILES)

clean:
	@echo "Cleaning up; deleting binaries..."
	@rm -rf $(BIN_DIR)
	@rm -rf $(APP_DATA)

wipe-data:
	@echo "Wiping app data..."
	@rm -rf $(APP_DATA)

lint:
	@golangci-lint run -c ./golangci.yml ./...

test:
	@go test ./... -v --cover

test-report:
	@go test ./... -v --cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out
