# Stage 1: Build stage
FROM golang:1.25.1-alpine AS build

ARG SERVICE_NAME

# Set environment variables
ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy all necessary files
COPY ./pkg ./pkg
COPY ./shared ./shared
COPY ./${SERVICE_NAME} ./${SERVICE_NAME}

RUN go work init ./pkg ./shared ./${SERVICE_NAME}

WORKDIR /app/${SERVICE_NAME}

# Install dependencies
RUN go mod tidy

# Build the application binary
RUN go build -o app cmd/main.go

# Stage 2: Create lightweight production image
FROM alpine:3.20.1 AS prod

ARG SERVICE_NAME

# Install dependencies
RUN apk add --no-cache tzdata

# Add application user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app

# Copy build binary from build stage
COPY --from=build /app/${SERVICE_NAME}/app .

# Change user
USER appuser

# Command to execute
CMD ["./app"]