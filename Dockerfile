# Stage 1: Build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/api

# Stage 2: Runtime enxuto
FROM alpine:3.20
WORKDIR /root/
COPY --from=builder /api .
EXPOSE 8080
CMD ["./api"]
