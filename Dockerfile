# Etapa 1: build
FROM golang:1.23.2 AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o rate-limiter ./cmd/server

# Etapa 2: imagem final
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/rate-limiter .
COPY .env .

EXPOSE 8080
CMD ["./rate-limiter"]