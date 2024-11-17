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

func AddServicesWithTransaction(w http.ResponseWriter, r *http.Request) {
	// Structure pour la réponse
	type Response struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Data    *models.Athlete `json:"data,omitempty"`
	}

	// Décoder le corps de la requête
	var body models.Athlete
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logs.Errorf("Erreur lors du décodage du JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Format de données invalide",
		})
		return
	}

	// Créer un contexte avec timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Démarrer la transaction
	tx, err := core.MysqlDb.BeginTx(ctx, nil)
	if err != nil {
		logs.Errorf("Erreur lors du démarrage de la transaction: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Erreur interne du serveur",
		})
		return
	}

	// Préparer la requête sans inclure l'ID car il est généré par la base de données
	query := `INSERT INTO athletes (nom, prenom, age) VALUES ($1, $2, $3) RETURNING id, nom, prenom, age`
	var insertedAthlete models.Athlete
	// Insérer l'athlète sans spécifier l'ID, qui sera généré par la base de données
	err = tx.QueryRowContext(ctx, query, body.Nom, body.Prenom, body.Age).Scan(&insertedAthlete.ID, &insertedAthlete.Nom, &insertedAthlete.Prenom, &insertedAthlete.Age)

	if err != nil {
		tx.Rollback()
		logs.Errorf("Erreur lors de l'insertion: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Erreur lors de l'ajout de l'athlète",
		})
		return
	}

	// Commit de la transaction
	if err := tx.Commit(); err != nil {
		logs.Errorf("Erreur lors du commit: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Erreur lors de la finalisation de l'ajout",
		})
		return
	}

	// Répondre avec succès
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Athlète ajouté avec succès",
		Data:    &insertedAthlete,
	})
}
