FROM golang:1.25-alpine

WORKDIR /app

# Установка зависимостей
RUN apk add --no-cache git

# Копируем файлы модулей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Устанавливаем инструмент для миграций
RUN wget https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz \
    && tar -xzf migrate.linux-amd64.tar.gz \
    && chmod +x migrate \
    && mv migrate /usr/local/bin/ \
    && rm migrate.linux-amd64.tar.gz

# Копируем Docker конфиг (если создали отдельный файл)
# COPY config.docker.toml config.toml

# Собираем приложение
RUN go build -o main ./cmd/apiserver/main.go

EXPOSE 5500

CMD ["./main"]