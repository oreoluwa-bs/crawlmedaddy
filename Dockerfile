# Step 1: build the binary
FROM golang:1.19-alpine as builder
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Step 2: build a small image
FROM alpine:latest
WORKDIR /root/

# Copy the compiled binary from the builder image
COPY --from=builder /app/main .

# Command to run the binary
CMD ["./crawlmedaddy/main.go"]
