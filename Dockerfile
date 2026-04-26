# Multi-stage build for minimal final image
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /stock-service ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /stock-service .
EXPOSE 3000
ENTRYPOINT ["/app/stock-service"]
