# Variables
APP_NAME=LibraChess

# Default target
all: build

# Build the Go binary
build:
	go build -o $(APP_NAME) main.go

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
	./dist/cutechess-cli/cutechess-cli \
		-engine name=PullLibra cmd=./libra \
		-engine name=MainLibra cmd=./libramain \
		-each proto=uci tc=180+2 \
		-games 100 \
		-concurrency 10 \
		-openings file=./books/chess.epd format=epd order=random plies=8 \
		-ratinginterval 10 \
		-draw movenumber=40 movecount=6 score=10 \
		-rounds 1