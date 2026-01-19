package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"todo/pkg/date"
	"todo/pkg/db"
)

func writeJSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

func jsonError(w http.ResponseWriter, msg string, statusCode int) {
	writeJSON(w, map[string]string{"error": msg}, statusCode)
}

func checkDate(task *db.Task) error {
	currentTime := time.Now()

	if task.Date == "" {
		task.Date = currentTime.Format(dateFormat)
	}

	_, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	var nextDate string
	if task.Repeat != "" {
		nextDate, err = date.NextDate(currentTime, task.Date, task.Repeat)
		if err != nil {
			return err
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
	if r.Method != "POST" {
		jsonError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newTask db.Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newTask); err != nil {
		jsonError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if newTask.Title == "" {
		jsonError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := checkDate(&newTask); err != nil {
		jsonError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	id, err := db.AddTask(&newTask)
	if err != nil {
		jsonError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"id": id}, http.StatusCreated)
}
