package api

import "net/http"

// Инициализация API
func Init() {
	// Вход в планировщик.
	http.HandleFunc("POST /api/signin", signinHandler)

	// Обработчик задачи.
	http.HandleFunc("POST /api/task", auth(addTaskHandler))

	// Обработчик выполнения задачи.
	http.HandleFunc("GET /api/task", auth(getTaskHandler))

	// Обработчик обновления задачи.
	http.HandleFunc("PUT /api/task", auth(updateTaskHandler))

	// Обработчик удаления задачи.
	http.HandleFunc("DELETE /api/task", auth(deleteTaskHandler))

	// Обработчик списка ближайших задач.
	http.HandleFunc("/api/task", auth(taskHandler))

	// Обработчик списка ближайших задач.
	http.HandleFunc("POST /api/task/done", auth(doneTaskHandler))

	// Обработчик списка ближайших задач.
	http.HandleFunc("GET /api/nextdate", nextDateHandler)

	// Обработчик списка ближайших задач.
	http.HandleFunc("GET /api/tasks", auth(tasksHandler))

	// Фронтенд.
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// Старый код, который был закомментирован. Он использовал другой способ регистрации обработчиков.
	// вход в планировщик
	//	http.HandleFunc("/api/signin", signinHandler)
	// обработчик задачи
	//	http.HandleFunc("/api/task", auth(taskHandler))
	// обработчик выполнения задачи
	//	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
	// обработчик следующей даты
	//	http.HandleFunc("/api/nextdate", nextDateHandler)
	// обработчик списка ближайших задач
	//	http.HandleFunc("/api/tasks", auth(tasksHandler))
	// фронтенд
	//	http.Handle("/", http.FileServer(http.Dir("./web")))

}

// обработчик неподдерживаемого метода для /api/task
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
