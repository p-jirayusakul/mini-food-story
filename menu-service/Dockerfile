# Stage 1: Build stage
FROM golang:1.24.2-alpine AS build

# Set environment variables
ENV CGO_ENABLED=0 \
    GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy Go workspace & shared libraries (pkg)
COPY ./go.work ./go.work.sum ./
COPY ./pkg ./pkg
COPY ./shared ./shared

# Copy directories ของทุก service ที่ถูกใช้งานใน workspace
COPY ./table-service ./table-service
COPY ./menu-service ./menu-service
COPY ./order-service ./order-service
COPY ./kitchen-service ./kitchen-service
COPY ./payment-service ./payment-service

# Set working directory to menu-service
WORKDIR /app/menu-service

# Install dependencies
RUN go mod tidy

# Build the application binary
RUN go build -o app cmd/main.go

# Stage 2: Create lightweight production image
FROM alpine:3.20.1 AS prod

RUN apk add --no-cache tzdata

# Set working directory in the production container
WORKDIR /app

# Copy compiled binary from build stage
COPY --from=build /app/menu-service/app .

ENV TZ=Asia/Bangkok

# Expose the port
EXPOSE 8081

# Set the default command to run the app
CMD ["./app"]