FROM golang:1.24-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o worker ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/worker .
COPY --from=builder /app/internal/config/config.yml /app/config.yml

EXPOSE 8080
ENTRYPOINT ["./worker"]