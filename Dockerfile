# сборка на alpine
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler .
# папка для базы, её потом можно заменить volume
RUN mkdir -p /data

# итоговый образ — пустой scratch
FROM scratch

WORKDIR /app

COPY --from=builder /app/scheduler /app/scheduler
# web берём из builder, а не с диска хоста
COPY --from=builder /app/web /app/web
COPY --from=builder /data /data

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db
# ENV TODO_PASSWORD=12345

# Задаем через переменную окружения, чтобы можно было менять при запуске контейнера 
# EXPOSE 7540
# база лежит в /data, том подключается при запуске
VOLUME /data

CMD ["/app/scheduler"]
