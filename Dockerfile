# Stage 1: Builder
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source
COPY . .

# Build the app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Stage 2: Final image
FROM alpine:latest

WORKDIR /app

# Add certs and migrate
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

# Expose the actual app port
EXPOSE 8081

CMD ["./main"]
