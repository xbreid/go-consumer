# Start from the official golang image
FROM golang:1.20.6-alpine

# Install git
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Create a directory in the container where the source code will reside
WORKDIR /app

# Add the source code into the container
ADD . /app

# Fetch the application dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 go build -o main ./cmd

RUN chmod +x /app/main

# Run the application binary.
CMD ["/app/main"]