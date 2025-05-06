FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install migrate tool
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Create a minimal image for running the application
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS requests and the migrate tool
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/images ./images
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]