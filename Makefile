# Variables
APP_NAME=libra-chess

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