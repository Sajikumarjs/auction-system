# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port on which your application runs
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
