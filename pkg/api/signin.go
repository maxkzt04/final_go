package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	// читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// возвращаем ошибку если не удалось прочитать тело запроса
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	// запрос от клиента
	var req map[string]string
	// парсим тело запроса
	err = json.Unmarshal(body, &req)
	// возвращаем ошибку если не удалось парсить тело запроса
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// пароль из переменной окружения
	pass := os.Getenv("TODO_PASSWORD")
	if req["password"] != pass || len(pass) == 0 {
		writeJson(w, http.StatusUnauthorized, map[string]string{"error": "Неверный пароль"})
		return
	}

	// создаем токен
	token, err := createToken(pass)
	if err != nil {
		// возвращаем ошибку если не удалось создать токен
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// возвращаем токен
	writeJson(w, http.StatusOK, map[string]string{"token": token})
}
