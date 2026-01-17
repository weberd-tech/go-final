package api

import (
	"net/http"

	"todo/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	allTasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if allTasks == nil {
		allTasks = []*db.Task{}
	}

	response := TasksResp{
		Tasks: allTasks,
	}
	writeJSON(w, response)
}
