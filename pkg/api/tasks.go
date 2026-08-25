package api

import (
	"net/http"

	"final.go/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// получаем параметр search
	search := r.URL.Query().Get("search")
	// получаем список задач
	tasks, err := db.Tasks(50, search)
	// возвращаем ошибку если не удалось получить список задач
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	// возвращаем список задач
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
