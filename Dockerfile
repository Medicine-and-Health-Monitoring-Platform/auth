# Stage 1: Build stage
FROM golang:1.22.5 AS builder

WORKDIR /app

# Copy and download dependencies
# COPY go.mod .
# COPY go.sum .

# Copy the rest of the application
COPY . .
RUN go mod download

# Optionally copy the .env file if it's needed
COPY .env .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -C ./cmd -a -installsuffix cgo -o ./../myapp .

# Stage 2: Final stage
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/myapp .


# Optionally copy the .env file if it's needed
COPY --from=builder /app/.env .

# Expose port 50052
EXPOSE 8081
EXPOSE 50051

# Command to run the executable
CMD ["./myapp"]