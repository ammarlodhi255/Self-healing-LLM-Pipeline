# Start with a lightweight Go image
FROM golang:1.19-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal image to run the application
FROM alpine:latest

# Set up a working directory in the final image
WORKDIR /app

# Install curl to help with debugging
RUN apk add --no-cache curl

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the static files and templates for the frontend
COPY --from=builder /app/static ./static
COPY --from=builder /app/templates ./templates

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the Go application
CMD ["./main"]