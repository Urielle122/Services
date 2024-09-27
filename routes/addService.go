package routes

import (
	"encoding/json"
	"net/http"
	"services/core"
	"services/models"
)

func AddServices(w http.ResponseWriter, r *http.Request) {
	body := models.ServiceModels{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO services (id, content, action, previous_content, title, type_action, type_action_name) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := core.MysqlDb.Exec(query, body.Id, body.Content, body.Action, body.PreviousContent, body.Title, body.TypeAction, body.TypeActionName)
	if err != nil {
		http.Error(w, "Failed to add service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
