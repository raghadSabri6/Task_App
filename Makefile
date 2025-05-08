.PHONY: build run test clean migrate-up migrate-down

# Build the application
build:
	go build -o bin/api cmd/api/main.go

# Run the application
run:
	go run cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run database migrations up
migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

# Run database migrations down
migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down