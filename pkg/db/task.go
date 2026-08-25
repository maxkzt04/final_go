package db

import (
	"database/sql"
	"fmt"
	"time"
)

// Задача из таблицы scheduler

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Добавление задачи в базу данных
func AddTask(task *Task) (int64, error) {

	// переменная для хранения ID задачи
	var id int64
	// запрос на добавление задачи в базу данных
	// date, title, comment, repeat - поля таблицы scheduler
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	// выполняем запрос на добавление задачи в базу данных
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)

	// если запрос выполнен успешно, получаем ID задачи
	if err == nil {
		// получаем ID задачи

		id, err = res.LastInsertId()
		// возвращаем ID задачи и ошибку
	}
	// возвращаем ID задачи и ошибку
	return id, err
	// если запрос выполнен не успешно, возвращаем 0 и ошибку
}

// получение списка ближайших задач
func Tasks(limit int, search string) ([]*Task, error) {

	var rows *sql.Rows
	var err error

	// если строка поиска пустая, возвращаем ближайшие задачи
	if search == "" {
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = db.Query(query, limit)
	} else {
		// проверяем, не дата ли это в формате 02.01.2006
		t, parseErr := time.Parse("02.01.2006", search)
		if parseErr == nil {
			// ищем задачи на конкретную дату
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`
			rows, err = db.Query(query, t.Format("20060102"), limit)
		} else {
			// поиск по заголовку и комментарию
			query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			search = "%" + search + "%"
			rows, err = db.Query(query, search, search, limit)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// пустой слайс, чтобы в JSON было {"tasks":[]}, а не {"tasks":null}
	tasks := make([]*Task, 0)

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

// получение задачи по id
func GetTask(id string) (*Task, error) {
	var task Task
	// нужно получить только одну запись, используем QueryRow
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// обновление задачи
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	// метод RowsAffected() возвращает количество записей к которым
	// была применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

// удаление задачи
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// обновление только даты задачи
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
