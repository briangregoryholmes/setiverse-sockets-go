# Use the official Go image as the base image for the build stage
FROM public.ecr.aws/docker/library/golang:latest AS build

# Set the working directory
WORKDIR /app

# Copy the necessary files to the working directory
COPY main.go go.mod go.sum ./

# Download and install the Go modules
RUN go mod download

# Build the Go application
RUN go build -o server .

# Use the official Alpine image as the base image for the final stage
FROM public.ecr.aws/docker/library/alpine:latest

# Create a non-root user and set ownership of the app directory
RUN adduser -D -g '' appuser
RUN mkdir /app && chown appuser:appuser /app

# Set the working directory
WORKDIR /app

# Copy the built binary from the build stage to the final stage
COPY --from=build --chown=appuser:appuser /app/server /app/server

# Ensure the built binary is executable
RUN chmod +x server

# Use the non-root user for running the container
USER appuser

# Expose the port that the application listens on
EXPOSE 4000

# Set the command to run when the container starts
CMD ["./server"]
