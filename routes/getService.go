package routes

import (
	"encoding/json"
	"net/http"
	"services/core"
	"services/models"
)

func GetServices(w http.ResponseWriter, r *http.Request) {
	body := models.ServiceModels{}
	val := len(body.PreviousContent)

	if val == 0{
		rows, err := core.MysqlDb.Query("SELECT id, title FROM services WHERE previous_content IS NULL")
		if err != nil {
			http.Error(w, "Failed to retrieve services", http.StatusInternalServerError)
		}
		defer rows.Close()

		var services []models.ServiceModels

	// Parcourir les résultats
		for rows.Next() {
			var service models.ServiceModels
			if err := rows.Scan(&service.Id, &service.Title); err != nil {
				http.Error(w, "Error scanning data", http.StatusInternalServerError)
				return
			}
			services = append(services, service)
			json.NewEncoder(w).Encode(services)
	}
}else {
    // Cas où val est un UUID
    rows, err := core.MysqlDb.Query("SELECT id, title FROM services WHERE previous_content = ?", body.PreviousContent)
    if err != nil {
        http.Error(w, "Failed to retrieve services", http.StatusInternalServerError)
        return
    }
    defer rows.Close()
    
    var services []models.ServiceModels
    for rows.Next() {
        var service models.ServiceModels
        if err := rows.Scan(&service.Id, &service.Title); err != nil {
            http.Error(w, "Error scanning data", http.StatusInternalServerError)
            return
        }
        services = append(services, service)
		json.NewEncoder(w).Encode(services)
    }
}
	w.Header().Set("Content-Type", "application/json")
	
}
