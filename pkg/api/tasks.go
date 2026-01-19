package api

import (
	"net/http"

	"todo/pkg/db"
)

const limitTasts = 50

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	allTasks, err := db.Tasks(limitTasts)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	if allTasks == nil {
		allTasks = []*db.Task{}
	}

	response := TasksResp{
		Tasks: allTasks,
	}
	writeJSON(w, response, http.StatusOK)
}
