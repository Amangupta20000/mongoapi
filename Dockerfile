# Use the official Golang image to create a build artifact.
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app. Ensure you set the output binary name to main
RUN go build -o main .

# Start a new stage from scratch using a minimal base image
FROM gcr.io/distroless/base

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Set the working directory in the final image
WORKDIR /

# Command to run the executable
CMD ["./main"]