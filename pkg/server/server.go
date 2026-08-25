package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"final.go/pkg/api"
)

// порт по умолчанию
const defaultPort = 7540

// Запуск сервера
// возвращает ошибку, если сервер не запустился
func Run() error {
	// Инициализация API
	api.Init()

	port := defaultPort
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		// парсим порт из переменной окружения
		if p, err := strconv.Atoi(envPort); err == nil {
			// устанавливаем порт из переменной окружения
			port = p
		}
	}
	// Запуск сервера
	// формируем адрес сервера
	addr := fmt.Sprintf(":%d", port)
	// выводим информацию о запуске сервера
	log.Printf("Server is running on %s", addr)
	// Возвращаем ошибку, если сервер не запустился
	return http.ListenAndServe(addr, nil)
}
