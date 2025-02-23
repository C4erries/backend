# Используем базовый образ Go
FROM golang:1.22.4-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Собираем приложение
RUN go build -o backend .

# alpine -- легковесный образ linux для контейнера
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем собранное приложение
COPY --from=builder /app/backend .

# Открываем порт для HTTP-запросов
EXPOSE ${APP_PORT}

# Запускаем приложение
CMD ["./backend"]