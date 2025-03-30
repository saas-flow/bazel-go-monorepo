# syntax=docker/dockerfile:1

# Base image for building the application
FROM golang:1.22 AS build-stage

WORKDIR /app

# Cache dependencies to speed up builds
COPY go.mod go.sum ./
RUN go mod download

# Copy source code separately to leverage Docker caching
COPY . .

# Build binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /docker-gs-ping

# Run tests in a separate stage to avoid bloating the final image
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Slim and secure runtime image
FROM gcr.io/distroless/static-debian11:nonroot AS build-release-stage

WORKDIR /

# Copy only the final binary from the build stage
COPY --from=build-stage /docker-gs-ping /docker-gs-ping

# Use non-root user for security
USER nonroot:nonroot

# Run the application
ENTRYPOINT ["/docker-gs-ping"]
