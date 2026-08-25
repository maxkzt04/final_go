package api

import (
	"net/http"
	"time"

	"final.go/pkg/db"
)

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	// выполнение задачи только методом POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// получаем идентификатор задачи
	id := r.URL.Query().Get("id")
	// проверяем идентификатор задачи
	if id == "" {
		writeJson(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// получаем задачу из базы данных
	task, err := db.GetTask(id)
	if err != nil {
		// возвращаем ошибку если задача не найдена
		writeJson(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	// если повторения нет — просто удаляем
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, map[string]string{"error": err.Error()})
			return
		}
		// возвращаем успешный ответ
		writeJson(w, map[string]string{})
		return
	}

	// периодическая задача — считаем следующую дату
	next, err := NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	// обновляем дату задачи
	err = db.UpdateDate(next, id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	// возвращаем успешный ответ
	writeJson(w, map[string]string{})
}
