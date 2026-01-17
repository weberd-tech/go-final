package api

import (
	"encoding/json"
	"net/http"

	"todo/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, foundTask)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var taskToUpdate db.Task

	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&taskToUpdate)
	if decodeErr != nil {
		writeJSON(w, map[string]string{"error": decodeErr.Error()})
		return
	}

	if taskToUpdate.ID == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	if taskToUpdate.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	checkErr := checkDate(&taskToUpdate)
	if checkErr != nil {
		writeJSON(w, map[string]string{"error": checkErr.Error()})
		return
	}

	updateErr := db.UpdateTask(&taskToUpdate)
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

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskId := r.FormValue("id")
	if taskId == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	deleteErr := db.DeleteTask(taskId)
	if deleteErr != nil {
		if deleteErr.Error() == "incorrect id for deleting task" {
			writeJSON(w, map[string]string{"error": "Задача не найдена"})
		} else {
			writeJSON(w, map[string]string{"error": deleteErr.Error()})
		}
		return
	}

	writeJSON(w, map[string]interface{}{})
}
