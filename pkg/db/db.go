package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Схема базы данных

const schema = `

CREATE TABLE scheduler (

    id INTEGER PRIMARY KEY AUTOINCREMENT,

    date CHAR(8) NOT NULL DEFAULT "",

    title VARCHAR(256) NOT NULL DEFAULT "",

    comment TEXT NOT NULL DEFAULT "",

    repeat VARCHAR(128) NOT NULL DEFAULT ""

);

CREATE INDEX idx_scheduler_date ON scheduler (date);

`

// Инициализация базы данных

// переменная для хранения соединения с базой данных
var db *sql.DB

// Инициализация базы данных
func Init() error {
	// получаем имя файла базы данных из переменной окружения
	dbFile := "scheduler.db"
	// если переменная окружения TODO_DBFILE не пустая, устанавливаем имя файла базы данных из переменной окружения
	if envFile := os.Getenv("TODO_DBFILE"); envFile != "" {
		dbFile = envFile
	}

	log.Println("файл базы:", dbFile)

	// в scratch папки /data нет — создаём, если нужно
	dir := filepath.Dir(dbFile)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// проверяем, существует ли файл базы данных
	_, err := os.Stat(dbFile)
	// если файл не существует, устанавливаем флаг install в true
	// иначе устанавливаем флаг install в false
	var install bool
	// если файл не существует, устанавливаем флаг install в true
	if err != nil {
		install = true
	}

	// открываем соединение с базой данных
	db, err = sql.Open("sqlite", dbFile)
	// возвращаем ошибку если не удалось открыть соединение

	if err != nil {
		return err
	}

	if install {
		// выполняем схему базы данных
		if _, err = db.Exec(schema); err != nil {
			// возвращаем ошибку если не удалось выполнить схему базы данных
			return err
		}
	}

	return nil
	// возвращаем nil если все прошло успешно
}
