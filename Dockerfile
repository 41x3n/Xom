# Start from the official Golang image
FROM golang:1.20-alpine

# Install ffmpeg and any other dependencies
RUN apk update && apk add --no-cache ffmpeg

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o /xom cmd/xom.go

# Specify the command to run the application
CMD ["/xom"]
