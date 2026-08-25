package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// afterNow возвращает true, если date больше now (без учёта времени)
func afterNow(date, now time.Time) bool {
	// форматируем дату в формате YYYYMMDD
	d1 := date.Format("20060102")
	// форматируем текущую дату в формате YYYYMMDD
	d2 := now.Format("20060102")
	return d1 > d2
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// пустое правило — ошибка
	if repeat == "" {
		return "", errors.New("не указано правило повторения")
	}
	// парсим дату в формате YYYYMMDD
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата начала: %s", dstart)
	}

	// разбиваем правило на части
	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("не указано правило повторения")
	}
	if parts[0] == "d" {
		// проверяем формат
		if len(parts) != 2 {
			return "", fmt.Errorf("неверный формат: %s", repeat)
		}
		// парсим интервал в днях
		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("некорректный интервал в днях: %s", parts[1])
		}
		// проверяем интервал в днях
		if days < 1 || days > 400 {
			return "", fmt.Errorf(
				"превышен максимально допустимый интервал: %d",
				days,
			)
		}

		// цикл для получения следующей даты
		for {
			// добавляем интервал в днях
			date = date.AddDate(0, 0, days)
			// проверяем, что дата больше текущей даты
			if afterNow(date, now) {
				break
			}
		}

	} else if parts[0] == "y" {
		// цикл для получения следующей даты
		for {
			date = date.AddDate(1, 0, 0)
			// проверяем, что дата больше текущей даты
			if afterNow(date, now) {
				break
			}
		}

	} else if parts[0] == "w" {
		// проверяем формат
		if len(parts) != 2 {
			return "", errors.New("неверный формат")
		}
		// массив для хранения дней недели
		var weekdays [8]bool
		wdays := strings.Split(parts[1], ",")

		// цикл для получения следующей даты
		for i := 0; i < len(wdays); i++ {
			// парсим день недели
			n, err := strconv.Atoi(wdays[i])
			if err != nil {
				return "", fmt.Errorf("некорректный день недели: %s", wdays[i])
			}
			// проверяем день недели
			if n < 1 || n > 7 {
				return "", fmt.Errorf("недопустимое значение: %d", n)
			}

			// устанавливаем флаг для дня недели
			weekdays[n] = true
		}

		// цикл для получения следующей даты
		for {
			// добавляем 1 день
			date = date.AddDate(0, 0, 1)
			// проверяем, что дата больше текущей даты
			if afterNow(date, now) {
				wd := int(date.Weekday())

				if wd == 0 {
					wd = 7
				}

				if weekdays[wd] {
					break
				}
			}
		}

	} else if parts[0] == "m" {
		// дни месяца, опционально месяцы
		if len(parts) < 2 || len(parts) > 3 {
			return "", fmt.Errorf("неверный формат: %s", repeat)
		}

		var day [32]bool
		var month [13]bool
		// флаг для последнего дня месяца
		minus1 := false
		// флаг для предпоследнего дня месяца
		minus2 := false

		// по умолчанию подходят все месяцы
		for i := 1; i <= 12; i++ {
			month[i] = true
		}
		// разбиваем список дней на части
		daysList := strings.Split(parts[1], ",")
		// цикл для получения следующей даты
		for i := 0; i < len(daysList); i++ {
			// парсим день месяца
			n, err := strconv.Atoi(daysList[i])
			if err != nil {
				return "", fmt.Errorf("некорректный день месяца: %s", daysList[i])
			}
			if n == -1 {
				minus1 = true
			} else if n == -2 {
				minus2 = true
			} else if n >= 1 && n <= 31 {
				// устанавливаем флаг для дня месяца
				day[n] = true
			} else {
				return "", fmt.Errorf("недопустимый день месяца: %d", n)
			}
		}

		// если указаны месяцы
		if len(parts) == 3 {
			// цикл для получения следующей даты
			for i := 1; i <= 12; i++ {
				month[i] = false
			}
			// разбиваем список месяцев на части
			monthsList := strings.Split(parts[2], ",")
			// цикл для получения следующей даты
			for i := 0; i < len(monthsList); i++ {
				n, err := strconv.Atoi(monthsList[i])
				if err != nil {
					return "", fmt.Errorf("некорректный месяц: %s", monthsList[i])
				}
				// проверяем, что месяц валиден
				if n < 1 || n > 12 {
					return "", fmt.Errorf("недопустимый месяц: %d", n)
				}
				month[n] = true
			}
		}

		// проверяем, что правило вообще будет выполнимо
		// например, m 30 2 или m 31 4 никогда не наступят согласно правилу
		daysInMonth := []int{0, 31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
		possible := minus1 || minus2
		for d := 1; d <= 31 && !possible; d++ {
			if !day[d] {
				continue
			}
			for m := 1; m <= 12; m++ {
				if month[m] && d <= daysInMonth[m] {
					possible = true
					break
				}
			}
		}
		if !possible {
			return "", fmt.Errorf("невозможно найти дату для правила: %s", repeat)
		}

		// счётчик, чтобы цикл не был бесконечным
		steps := 0
		for {
			steps++
			if steps > 366*5 {
				return "", fmt.Errorf("не удалось найти следующую дату")
			}
			// добавляем 1 день
			date = date.AddDate(0, 0, 1)
			// проверяем, что дата больше текущей даты
			if afterNow(date, now) {
				m := int(date.Month())
				d := date.Day()
				// последний день текущего месяца
				last := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
				// флаг для проверки, что дата валидна
				ok := false
				// проверяем, что месяц валиден
				if month[m] {
					// проверяем, что день валиден
					if day[d] {
						// проверяем, что день валиден
						ok = true
					}
					if minus1 && d == last {
						// проверяем, что день валиден
						ok = true
					}
					if minus2 && d == last-1 {
						// проверяем, что день валиден
						ok = true
					}
				}
				if ok {
					// если дата валидна, выходим из цикла
					break
				}
			}
		}
	} else {
		// возвращаем ошибку если формат не валиден
		return "", fmt.Errorf("неверный формат: %s", repeat)
	}

	// форматируем дату в формате YYYYMMDD
	return date.Format("20060102"), nil
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// получаем параметр now
	nowStr := r.URL.Query().Get("now")
	// получаем параметр date
	date := r.URL.Query().Get("date")
	// получаем параметр repeat
	repeat := r.URL.Query().Get("repeat")
	// инициализируем переменные
	var now time.Time
	// переменная для хранения ошибки
	var err error
	// проверяем, что параметр now указан
	if nowStr == "" {
		// если не указан параметр now, устанавливаем текущую дату
		now = time.Now()
	} else {
		// парсим дату в формате YYYYMMDD
		now, err = time.Parse("20060102", nowStr)
		// возвращаем ошибку если не удалось парсить дату
		if err != nil {
			// возвращаем ошибку в формате JSON
			fmt.Fprint(w, err.Error())
			return
		}
	}
	// получаем следующую дату
	next, err := NextDate(now, date, repeat)
	if err != nil {
		// возвращаем ошибку в формате JSON
		fmt.Fprint(w, err.Error())
		return
	}
	// возвращаем результат в формате JSON
	fmt.Fprint(w, next)
}
