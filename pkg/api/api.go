package api

import (
	"net/http"
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("GET /api/task", getTaskHandler)
	http.HandleFunc("POST /api/task", addTaskHandler)
	http.HandleFunc("PUT /api/task", updateTaskHandler)
	http.HandleFunc("DELETE /api/task", deleteTaskHandler)
	http.HandleFunc("POST /api/task/done", onTaskDone)
	http.HandleFunc("/api/tasks", tasksHandler)
}
