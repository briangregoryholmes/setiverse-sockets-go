# Use the official Go image as the base image
FROM golang:1.17.5-alpine3.15

# Set the working directory
WORKDIR /app

# Copy the necessary files to the working directory
COPY main.go go.mod go.sum ./

# Download and install the Go modules
RUN go mod download

# Build the Go application
RUN go build -o server .

# Ensure the built binary is executable
RUN chmod +x server

# Expose the port that the application listens on
EXPOSE 4000

# Set the command to run when the container starts
CMD ["./server"]
