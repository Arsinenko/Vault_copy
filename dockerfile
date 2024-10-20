# Используем официальный образ Go
FROM golang:1.23 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы приложения
COPY . .

# Компилируем приложение
RUN go build -o myapp .

# Используем легковесный образ для запуска приложения
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates

# Копируем скомпилированное приложение из предыдущего этапа
COPY --from=builder /app/myapp /myapp

# Указываем команду для запуска приложения
CMD ["/myapp"]