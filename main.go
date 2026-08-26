package main

import (
	"log"
	"os"

	"final.go/pkg/db"
	"final.go/pkg/server"
	"github.com/joho/godotenv"
)

// Не забыть добавить в .env переменные окружения для базы данных и порта
func main() {
	// пишем логи в stdout — так их видно в docker logs
	log.SetOutput(os.Stdout)

	err := godotenv.Load()
	if err != nil {
		log.Println("не удалось загрузить .env")
	}

	// Инициализация базы данных
	if err := db.Init(); err != nil {
		log.Println("ошибка базы данных:", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("ошибка закрытия базы данных:", err)
		}
	}()

	// Запуск сервера
	if err := server.Run(); err != nil {
		log.Println("ошибка сервера:", err)
		// тут возник вопрос, а нужен ли тут return? ведь main() и так завершится после log.Println,
		// решил его оставить, return чтобы явно показать что выполнение программы прекращается
		// по заметки Ревьюера
		return
	}
}
