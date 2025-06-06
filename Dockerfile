# Сборочный этап
FROM golang:1.24.1 AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y git

# Загрузка зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка бинарника (важно указать корректный путь к main пакету)
RUN go build -o consumer ./cmd/app

# Финальный образ
FROM gcr.io/distroless/base-debian12

WORKDIR /app

# Добавление переменных окружения
COPY .env .env

# Копируем скомпилированный бинарник
COPY --from=builder /app/consumer /app/consumer

# Указываем команду запуска
CMD ["/app/consumer"]
