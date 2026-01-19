package api

import (
	"encoding/json"
	"net/http"

	"todo/pkg/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	writeJSON(w, foundTask, http.StatusOK)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var taskToUpdate db.Task

	decoder := json.NewDecoder(r.Body)
	decodeErr := decoder.Decode(&taskToUpdate)
	if decodeErr != nil {
		writeJSON(w, map[string]string{"error": decodeErr.Error()}, http.StatusBadRequest)
		return
	}

	if taskToUpdate.ID == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	if taskToUpdate.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	checkErr := checkDate(&taskToUpdate)
	if checkErr != nil {
		writeJSON(w, map[string]string{"error": checkErr.Error()}, http.StatusUnprocessableEntity)
		return
	}

	updateErr := db.UpdateTask(&taskToUpdate)
	if updateErr != nil {
		if updateErr.Error() == "incorrect id for updating task" {
			writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]string{"error": updateErr.Error()}, http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskId := r.FormValue("id")
	if taskId == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	deleteErr := db.DeleteTask(taskId)
	if deleteErr != nil {
		if deleteErr.Error() == "incorrect id for deleting task" {
			writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		} else {
			writeJSON(w, map[string]string{"error": deleteErr.Error()}, http.StatusInternalServerError)
		}
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
