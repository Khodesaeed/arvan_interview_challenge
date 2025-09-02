# Stage 1: Build the Go application
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-app ./main.go

# Stage 2: Create the final production image
FROM alpine:3.19

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /go-app .

# Expose the port the app will listen on
EXPOSE 8080 5050

# Run the application
CMD ["./go-app"]
