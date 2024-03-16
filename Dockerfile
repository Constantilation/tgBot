# Используйте образ на основе Debian или другого дистрибутива Linux, который содержит необходимые инструменты для сборки C-программ
FROM golang:1.21rc3 as builder

# Установка зависимостей для CGO
RUN apt-get update && apt-get install -y default-libmysqlclient-dev gcc

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Убедитесь, что CGO_ENABLED по умолчанию включен
RUN CGO_ENABLED=1 GOOS=linux go build -o bot .

FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

WORKDIR /root/

COPY --from=builder /app/bot .
COPY --from=builder /app/.env .

CMD ["ls -la"]

CMD ["./bot"]
