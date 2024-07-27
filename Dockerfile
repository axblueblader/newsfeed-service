# Use a lightweight, minimal base image
FROM docker.io/golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the executable
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a smaller base image for production
FROM alpine:latest
RUN adduser -D appuser

# Copy the built executable
COPY --from=builder /app/main /app/
RUN chown -R appuser:appuser /app

USER appuser

# Expose the port
EXPOSE 8080
EXPOSE 5432

# Command to run the application
CMD ["./app/main"]
