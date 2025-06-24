# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o service-a main.go

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/service-a .
# gRRPC & Prometheus
EXPOSE 50051 9091 
CMD ["./service-a"]
