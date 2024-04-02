FROM golang:latest as builder
LABEL authors="@Zatraaz"

# Создание рабочий директории
RUN mkdir -p /app

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы проекта внутрь контейнера
COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build -o go_test ./cmd/app/main.go

# Второй этап: создание production образ
FROM ubuntu AS chemistry

WORKDIR /app

RUN apt-get update

COPY --from=builder /app/go_test ./
COPY ./ ./

CMD ["./go_test"]