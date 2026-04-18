# Этап 1: Сборка
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/main.go

# Этап 2: Минимальный образ для запуска
FROM alpine:latest

WORKDIR /app

# Устанавливаем tzdata для корректной работы с временем
RUN apk --no-cache add tzdata

# Копируем бинарник из этапа сборки
COPY --from=builder /app/server /app/server

# Копируем статические файлы (фронтенд)
COPY --from=builder /app/static /app/static

# Открываем порт
EXPOSE 8080

# Запускаем приложение
# Переменные окружения передаются через docker-compose.yml (environment:)
CMD ["/app/server"]
