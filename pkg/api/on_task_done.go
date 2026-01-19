package api

import (
	"net/http"
	"strconv"
	"time"

	"todo/pkg/date"
	"todo/pkg/db"
)

func onTaskDone(w http.ResponseWriter, r *http.Request) {
	taskId := r.FormValue("id")
	if taskId == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	foundTask, err := db.GetTask(taskId)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}

	if foundTask.Repeat == "" {
		deleteErr := db.DeleteTask(taskId)
		if deleteErr != nil {
			writeJSON(w, map[string]string{"error": deleteErr.Error()}, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]interface{}{}, http.StatusOK)
		return
	}

	now := time.Now()
	nextDate, nextErr := date.NextDate(now, foundTask.Date, foundTask.Repeat)
	if nextErr != nil {
		writeJSON(w, map[string]string{"error": nextErr.Error()}, http.StatusInternalServerError)
		return
	}

	taskIdInt, parseErr := strconv.ParseInt(taskId, 10, 64)
	if parseErr != nil {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	updateErr := db.UpdateDate(nextDate, taskIdInt)
	if updateErr != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
