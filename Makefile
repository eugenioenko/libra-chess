# Variables
APP_NAME=libra-chess

# Default target
all: build

# Build the Go binary
build:
	go build -o $(APP_NAME) main.go

# Build release binaries for Linux, macOS, and Windows
.PHONY: build-release
build-release:
	mkdir -p release
	# Linux
	GOOS=linux GOARCH=amd64 go build -o release/$(APP_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o release/$(APP_NAME)-linux-arm64 main.go
	GOOS=linux GOARCH=386 go build -o release/$(APP_NAME)-linux-386 main.go
	# macOS
	GOOS=darwin GOARCH=amd64 go build -o release/$(APP_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o release/$(APP_NAME)-darwin-arm64 main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -o release/$(APP_NAME)-windows-amd64.exe main.go
	GOOS=windows GOARCH=arm64 go build -o release/$(APP_NAME)-windows-arm64.exe main.go
	GOOS=windows GOARCH=386 go build -o release/$(APP_NAME)-windows-386.exe main.go

# Run the application
run:
	go run main.go

# Format code using gofmt
fmt:
	gofmt -w .

# Run tests
.PHONY: test
test:
	go test -p 1 -v ./...

# Run a linter (requires `golangci-lint`)
lint:
	golangci-lint run ./...

test-cutechess:
	make build
	./dist/cutechess-cli/cutechess-cli \
		-engine name=PullLibra cmd=./libra-chess \
		-engine name=MainLibra cmd=./libra-main \
		-each proto=uci tc=180+2 \
		-games 10 \
		-concurrency 10 \
		-openings file=./books/chess.epd format=epd order=random plies=8 \
		-ratinginterval 10 \
		-draw movenumber=40 movecount=6 score=10 \
		-rounds 1

test-stockfish:
	make build
	./dist/cutechess-cli/cutechess-cli \
		-engine name=PullLibra cmd=./libra-chess \
		-engine name=Stockfish cmd=./stockfish/stockfish-cli option.UCI_LimitStrength=true option.UCI_Elo=1320 \
		-each proto=uci tc=10+0.1 \
		-games 10 \
		-concurrency 10 \
		-openings file=./books/chess.epd format=epd order=random plies=8 \
		-ratinginterval 10 \
		-draw movenumber=40 movecount=6 score=10 \
		-rounds 1

test-debug:
	make build
	./dist/cutechess-cli/cutechess-cli \
		-engine name=PullLibra cmd=./libra-chess \
		-engine name=MainLibra cmd=./libra-main \
		-each proto=uci tc=180+2 \
		-games 1 \
		-concurrency 10 \
		-ratinginterval 1 \
		-draw movenumber=40 movecount=6 score=10 \
		-debug \
		-rounds 1