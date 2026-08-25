package api

import "net/http"

// Инициализация API
func Init() {
	// вход в планировщик
	http.HandleFunc("/api/signin", signinHandler)

	// обработчик задачи
	http.HandleFunc("/api/task", auth(taskHandler))

	// обработчик выполнения задачи
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))

	// обработчик следующей даты
	http.HandleFunc("/api/nextdate", nextDateHandler)

	// обработчик списка ближайших задач
	http.HandleFunc("/api/tasks", auth(tasksHandler))

	// фронтенд
	http.Handle("/", http.FileServer(http.Dir("./web")))
}

// обработчик задачи
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// добавление задачи
		addTaskHandler(w, r)
	case http.MethodGet:
		// получение задачи
		getTaskHandler(w, r)
	case http.MethodPut:
		// обновление задачи
		updateTaskHandler(w, r)
	case http.MethodDelete:
		// удаление задачи
		deleteTaskHandler(w, r)
	}
}
