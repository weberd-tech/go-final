package api

import (
	"encoding/json"
	"net/http"
	"time"

	"todo/pkg/date"
	"todo/pkg/db"
)

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encoder := json.NewEncoder(w)
	encoder.Encode(data)
}

func checkDate(task *db.Task) error {
	currentTime := time.Now()

	if task.Date == "" {
		task.Date = currentTime.Format(dateFormat)
	}

	_, parseErr := time.Parse(dateFormat, task.Date)
	if parseErr != nil {
		return parseErr
	}

	var nextDate string
	if task.Repeat != "" {
		nextDate, parseErr = date.NextDate(currentTime, task.Date, task.Repeat)
		if parseErr != nil {
			return parseErr
		}
	}

	today := currentTime.Format(dateFormat)
	if task.Date < today {
		if task.Repeat == "" {
			task.Date = today
		} else {
			task.Date = nextDate
		}
	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask db.Task

	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&newTask)
	if decodeErr != nil {
		writeJSON(w, map[string]string{"error": decodeErr.Error()})
		return
	}

	if newTask.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	checkErr := checkDate(&newTask)
	if checkErr != nil {
		writeJSON(w, map[string]string{"error": checkErr.Error()})
		return
	}

	taskId, addErr := db.AddTask(&newTask)
	if addErr != nil {
		writeJSON(w, map[string]string{"error": addErr.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": taskId})
}
