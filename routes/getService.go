package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"services/core"
	"services/logs"
	"services/models"
	"time"
)

func GetAllAthletes(w http.ResponseWriter, r *http.Request) {
	// Structure pour la réponse
	type Response struct {
		Success bool              `json:"success"`
		Message string           `json:"message"`
		Data    []models.Athlete `json:"data,omitempty"`
	}

	// Créer un contexte avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Préparer la requête
	query := `SELECT id, nom, prenom, age FROM athletes`
	
	// Exécuter la requête
	rows, err := core.MysqlDb.QueryContext(ctx, query)
	if err != nil {
		logs.Errorf("Erreur lors de la récupération des athlètes: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Erreur lors de la récupération des athlètes",
		})
		return
	}
	defer rows.Close()

	// Slice pour stocker les athlètes
	var athletes []models.Athlete

	// Parcourir les résultats
	for rows.Next() {
		var athlete models.Athlete
		if err := rows.Scan(&athlete.ID, &athlete.Nom, &athlete.Prenom, &athlete.Age); err != nil {
			logs.Errorf("Erreur lors du scan d'un athlète: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Erreur lors de la lecture des données",
			})
			return
		}
		athletes = append(athletes, athlete)
	}

	// Vérifier les erreurs de parcours
	if err = rows.Err(); err != nil {
		logs.Errorf("Erreur lors du parcours des résultats: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Erreur lors de la lecture des données",
		})
		return
	}

	// Répondre avec succès
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Liste des athlètes récupérée avec succès",
		Data:    athletes,
	})
}