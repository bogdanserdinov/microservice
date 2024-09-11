# Base stage for building the application
FROM golang:1.21.6-alpine AS builder

# Set the working directory
WORKDIR /usr/src/app

# Cache the go.mod and go.sum files to leverage caching of dependencies
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the source code (excluding the 'main.go' to benefit from caching)
COPY . .

# Copy the main.go separately (so changes in the code don't invalidate the dependency layer)
COPY ./cmd/main.go ./cmd/main.go

# Build the Go application
RUN GOOS=linux go build -o microservice ./cmd/main.go

# Production stage
FROM alpine:3 AS production

# Set the working directory
ENV HOME /usr/src/app
WORKDIR $HOME

# Copy the built binary from the builder stage
COPY --from=builder /usr/src/app/microservice ./microservice
