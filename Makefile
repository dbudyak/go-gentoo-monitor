.PHONY: build run clean docker-build docker-up docker-down test

build:
	go build -o monitor ./cmd/monitor

run: build
	./monitor

clean:
	rm -f monitor
	go clean

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

test:
	go test -v ./...

deps:
	go mod download
	go mod tidy
