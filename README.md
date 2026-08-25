# Планировщик задач

Веб-сервер на Go: задачи с повторением, поиском и паролем. База — SQLite (`scheduler.db`), фронтенд в папке `web`.

## Задания со звёздочкой

- полный расчёт следующей даты (`d`, `y`, `w`, `m`)
- поиск по тексту и по дате `02.01.2006`
- аутентификация (`TODO_PASSWORD` + JWT)
- Docker

## Требования

Go **1.26.1** или выше:

```
go version
```

## Локальный запуск

```
go mod download
go run main.go
```

Сайт: http://localhost:7540

Переменные читаются из `.env` (godotenv). Пример — `.env.example`:

```
TODO_DBFILE=scheduler.db
TODO_PORT=7540
TODO_PASSWORD=12345
```

| Переменная | По умолчанию | Смысл |
| --- | --- | --- |
| `TODO_PORT` | `7540` | порт сервера |
| `TODO_DBFILE` | `scheduler.db` | файл SQLite |
| `TODO_PASSWORD` | пусто | пароль; если нет — вход не нужен |

Если задан `TODO_PASSWORD`, вход: http://localhost:7540/login.html  
После смены `.env` сервер нужно **перезапустить** (`Ctrl+C`, снова `go run main.go`).

## Тесты

Сначала запусти сервер (`go run main.go`), потом тесты. Без сервера API-тесты не пройдут.

Проверить всё сразу:

```
go test -v -count=1 ./tests
```

или по шагам (сервер уже должен быть запущен):

```
go test -run ^TestDB$ ./tests
go test -run ^TestNextDate$ ./tests
go test -run ^TestAddTask$ ./tests
go test -run ^TestTasks$ ./tests
go test -run ^TestTask$ ./tests
go test -run ^TestEditTask$ ./tests
go test -run ^TestDone$ ./tests
go test -run ^TestDelTask$ ./tests
```

После `go run main.go` в корне появляется `scheduler.db` — можно открыть любой программой для SQLite и посмотреть таблицу `scheduler`. Добавление, поиск, редактирование, выполнение и удаление дополнительно проверяются в браузере: http://localhost:7540

Настройки в `tests/settings.go`:

- `Port = 7540` — порт сервера
- `DBFile = "../scheduler.db"` — путь к базе относительно `tests/` (или переменная `TODO_DBFILE`)
- `FullNextDate = true` — тест `TestNextDate` проверяет ещё правила `w` и `m`
- `Search = true` — тест `TestTasks` проверяет поиск
- `Token` — JWT, если включён пароль; иначе пустая строка ``

### Тесты без пароля

В `.env` закомментируй пароль:

```
# TODO_PASSWORD=12345
```

Перезапусти сервер. В `settings.go` оставь `Token = ```.  
Если ошибка `invalid character 'A'` — на порту 7540 ещё старый процесс с паролем. Убей его и запусти сервер заново:

```
netstat -ano | findstr :7540
taskkill /PID <PID> /F
go run main.go
```

### Тесты с паролем

1. В `.env` раскомментируй `TODO_PASSWORD=12345`, перезапусти сервер.
2. Получи токен.

**Linux / macOS:**

```
curl -s -X POST http://localhost:7540/api/signin \
  -H "Content-Type: application/json" \
  -d '{"password":"12345"}'
```

**cmd (Windows):**

```
echo {"password":"12345"}>signin.json
curl.exe -s -X POST http://localhost:7540/api/signin -H "Content-Type: application/json" --data-binary @signin.json
```

**PowerShell:**

```
Set-Content -Path signin.json -Value '{"password":"12345"}' -Encoding ascii
curl.exe -s -X POST http://localhost:7540/api/signin -H "Content-Type: application/json" --data-binary "@signin.json"
```

В ответе будет `{"token":"..."}`. Это значение вставь в `tests/settings.go`:

```
var Token = `сюда_весь_токен`
```

В cmd/PowerShell не оборачивай JSON в одинарные кавычки — тело запроса ломается. В Linux/macOS одинарные кавычки как раз нужны.

## Docker

Сборка на `golang:alpine`, финальный образ — `scratch`. Данные в `/data` (`VOLUME /data`). Если контейнер упал: `docker logs scheduler`.

Сборка:

```
docker build -t scheduler .
```

Если контейнер с таким именем уже есть:

```
docker rm -f scheduler
```

Запуск (база и `.env` с хоста):

```
docker run --name scheduler -p 7540:7540 -v ${PWD}:/data --env-file .env scheduler
```

На Windows путь с кириллицей иногда ломает bind-mount — тогда том Docker:

```
docker volume create data
docker run --name scheduler -p 7540:7540 -v data:/data --env-file .env scheduler
```

Сайт: http://localhost:7540

Пароль в образ при сборке не кладётся — только через `--env-file .env` или `-e TODO_PASSWORD=...` при запуске.

## Полезные команды

Если порт 7540 занят:

**Windows:**

```
netstat -ano | findstr :7540
```

Примерный вывод:

```
  TCP    0.0.0.0:7540           0.0.0.0:0              LISTENING       30644
```

Убить процесс (подставь свой PID):

```
taskkill /PID 30644 /F
```

**Linux / macOS:**

```
ss -lptn 'sport = :7540'
```

или:

```
lsof -i :7540
```

Убить процесс:

```
kill 30644
```

если не завершился:

```
kill -9 30644
```

Остальные команды:

```
go mod verify
```
Проверяет, что зависимости в `go.sum` не повреждены и совпадают с модулями.

```
go mod tidy
```
Добавляет недостающие зависимости в `go.mod`/`go.sum` и убирает лишние.

```
go run main.go
```
Собирает и сразу запускает сервер. Сайт: http://localhost:7540

```
gofmt -w .
```
Форматирует весь Go-код в проекте.

```
go build ./...
```
Проверяет, что проект собирается (бинарник в текущую папку не обязательно остаётся).

```
go test -v -count=1 ./tests
```
Запускает тесты из `tests` с подробным выводом. `-count=1` — без кэша. Сервер при этом должен быть уже запущен.

```
docker build -t scheduler .
```
Собирает Docker-образ с именем `scheduler` из текущего `Dockerfile`.

```
docker logs scheduler
```
Показывает логи контейнера `scheduler` (удобно, если контейнер упал).

```
docker rm -f scheduler
```
Останавливает и удаляет контейнер `scheduler`, если он уже есть (перед новым `docker run`).
