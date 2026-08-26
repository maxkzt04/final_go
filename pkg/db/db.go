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
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(256) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
`

// переменная для хранения соединения с базой данных
var db *sql.DB

// Инициализация базы данных
func Init() error {
	// получаем имя файла базы данных из переменной окружения
	dbFile := "scheduler.db"

	// если переменная окружения TODO_DBFILE не пустая,
	// устанавливаем имя файла базы данных из переменной окружения
	if envFile := os.Getenv("TODO_DBFILE"); envFile != "" {
		dbFile = envFile
	}

	log.Println("файл базы:", dbFile)

	// создаём каталог для базы данных, если его нет
	dir := filepath.Dir(dbFile)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// открываем существующую базу данных или создаём новую
	var err error
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// создаём таблицу и индекс, если их ещё нет
	if _, err = db.Exec(schema); err != nil {
		_ = db.Close()
		return err
	}

	return nil
}

// Закрытие базы данных
func Close() error {
	if db == nil {
		return nil
	}

	return db.Close()
}

// Получение соединения с базой данных
func GetDB() *sql.DB {
	return db
}
