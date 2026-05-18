# Use the latest Go version to satisfy go 1.25+
FROM golang:latest AS builder

WORKDIR /app

# Install dependencies first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o controlplane cmd/controlplane/main.go

# Use a lightweight final image
FROM alpine:latest
WORKDIR /app

# The Alpine image needs libc compatibility for standard Go binaries
# when they are built with 'latest' instead of alpine directly
RUN apk add --no-cache gcompat

# Copy the built binary and the .env file
COPY --from=builder /app/controlplane .
COPY --from=builder /app/.env .

# Expose the API port
EXPOSE 8080

# Run the app
CMD ["./controlplane"]
