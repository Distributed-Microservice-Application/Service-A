# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build with optimizations for smaller binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o service-a main.go

# Run stage - use distroless for even smaller size
FROM gcr.io/distroless/static-debian11:nonroot
WORKDIR /
COPY --from=builder /app/service-a .
# gRPC & Prometheus
EXPOSE 50051 9091 
USER nonroot:nonroot
ENTRYPOINT ["./service-a"]
