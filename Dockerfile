# Этап сборки
FROM golang:1.23.4-alpine AS builder

# Установим рабочую директорию для сборки
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загружаем все зависимости
RUN go mod tidy || { echo 'go mod tidy failed'; exit 1; }

# Копируем исходный код в контейнер
COPY . .

# Собираем приложение
RUN go build -o /app/cmd/server/server cmd/server/main.go

# Этап финального изображения
FROM alpine:latest

# Устанавливаем необходимые библиотеки (если нужно, например, для работы с сертификатами или правами)
RUN apk --no-cache add ca-certificates


# Копируем скомпилированный бинарник из этапа сборки
COPY --from=builder /app/cmd/server/server /app/server

# Определяем порт, на котором будет работать приложение
EXPOSE 8080

# Запускаем приложение
CMD ["/app/server"]
