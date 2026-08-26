package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"final.go/pkg/db"
)

// ответ в JSON для API
func writeJson(w http.ResponseWriter, status int, data any) {
	resp, err := json.Marshal(data)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Ошибка формирования ответа"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_, _ = w.Write(resp)
}

// проверка даты задачи
// возвращает ошибку если дата некорректна
func checkDate(task *db.Task) error {
	now := time.Now()

	// если поле date не указано, берём сегодняшнее число в формате YYYYMMDD
	if task.Date == "" {
		task.Date = now.Format(db.DateFormat)
	}
	// парсим дату в формате YYYYMMDD
	t, err := time.Parse(db.DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректная дата: %s", task.Date)
	}

	var next string
	// если правило повторения указано, проверяем его
	if len(task.Repeat) > 0 {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			// возвращаем ошибку если не удалось получить следующую дату
			return err
		}
	}

	// если сегодня (now) больше task.Date (t)
	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			// если правила повторения нет, то берём сегодняшнее число
			task.Date = now.Format(db.DateFormat)
		} else {
			// в противном случае, берём вычисленную ранее следующую дату
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// возвращаем ошибку если не удалось прочитать тело запроса
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// парсим тело запроса в структуру Task
	err = json.Unmarshal(body, &task)
	// возвращаем ошибку если не удалось парсить тело запроса
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// заголовок обязательный
	if len(task.Title) == 0 {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// проверяем дату задачи
	err = checkDate(&task)
	if err != nil {
		// возвращаем ошибку если дата некорректна
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		// возвращаем ошибку если не удалось добавить задачу
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// возвращаем ID задачи
	writeJson(w, http.StatusOK, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}

	writeJson(w, http.StatusOK, task)
}

// обновление задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	// читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	// парсим тело запроса в структуру Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	// проверяем идентификатор задачи
	if task.ID == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	// проверяем заголовок задачи
	if len(task.Title) == 0 {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}
	// проверяем дату задачи
	err = checkDate(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	// обновляем задачу
	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	// возвращаем успешный ответ
	writeJson(w, http.StatusOK, map[string]string{})
}

// удаление задачи
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// получаем идентификатор задачи
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}
	// удаляем задачу
	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	// возвращаем успешный ответ
	writeJson(w, http.StatusOK, map[string]string{})
}
