APP_NAME = linebot
BUILD_DIR = $(PWD)/build

clean:
	rm -rf ./build

critic:
	gocritic check ./...

critic.all:
	gocritic check -enableAll ./...
security:
	gosec ./...

lint:
	golangci-lint run ./...

test: clean critic security lint
	go clean -testcache && go test -v -timeout 30s -coverprofile=cover.out -cover ./tests/...
	go tool cover -func=cover.out

build: clean critic security lint
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) ./cmd/app

run: build
	$(BUILD_DIR)/$(APP_NAME)

generate:
	go build -o $(BUILD_DIR)/generate ./cmd/generate

