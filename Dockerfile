# Этап сборки
FROM golang:1.23-alpine AS builder

# Установка рабочей директории
WORKDIR /app

# Включаем кэширование для go build
#ENV GOCACHE=/root/.cache/go-build
#
# Оптимизация для размера бинарного файла
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Установка необходимых инструментов для сборки
RUN apk add --no-cache git make && \
    go install github.com/swaggo/swag/cmd/swag@v1.16.3

# Копирование файлов
COPY go.mod go.sum ./
COPY vendor/ ./vendor/
COPY . .

# Создание директории для swagger если её нет
RUN mkdir -p docs/swagger

# Генерация swagger документации
RUN swag init \
    --parseDependency \
    --parseInternal \
    --parseDepth 5 \
    -g internal/app/start/start.go \
    --output docs/swagger 

# Сборка приложения
RUN go build -mod=vendor -ldflags="-s -w" -o ./bin/app ./cmd/app/*

# Финальный этап
FROM scratch

# Копируем сертификаты, пользователей и собранный бинарник
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/bin/app /app/
COPY --from=builder /app/docs/swagger /app/docs/swagger

# Установка рабочей директории
WORKDIR /app

# Переключение на непривилегированного пользователя
USER nobody

# Открытие порта
EXPOSE 8080

# Команда запуска приложения
CMD ["/app/app"]