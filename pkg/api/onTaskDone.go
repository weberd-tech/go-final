package api

import (
	"net/http"
	"time"

	"todo/pkg/date"
	"todo/pkg/db"
)

func onTaskDone(w http.ResponseWriter, r *http.Request) {
	taskId := r.FormValue("id")
	if taskId == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	foundTask, err := db.GetTask(taskId)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	if foundTask.Repeat == "" {
		deleteErr := db.DeleteTask(taskId)
		if deleteErr != nil {
			writeJSON(w, map[string]string{"error": deleteErr.Error()})
			return
		}
		writeJSON(w, map[string]interface{}{})
		return
	}

	now := time.Now()
	nextDate, nextErr := date.NextDate(now, foundTask.Date, foundTask.Repeat)
	if nextErr != nil {
		writeJSON(w, map[string]string{"error": nextErr.Error()})
		return
	}

	updateErr := db.UpdateDate(nextDate, taskId)
	if updateErr != nil {
		if updateErr.Error() == "incorrect id for updating task" {
			writeJSON(w, map[string]string{"error": "Задача не найдена"})
		} else {
			writeJSON(w, map[string]string{"error": updateErr.Error()})
		}
		return
	}

	writeJSON(w, map[string]interface{}{})
}
