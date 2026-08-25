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
		log.Fatal(err)
	}

	// Запуск сервера
	if err := server.Run(); err != nil {
		log.Println("ошибка сервера:", err)
		log.Fatal(err)
	}
}
