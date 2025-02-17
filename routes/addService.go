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

func removeFileExtension(fileName string) string {
	// Trouve la dernière occurrence du point "."
	lastDotIndex := strings.LastIndex(fileName, ".")                                                                                                                                                                                          
	if lastDotIndex == -1 {
		// Si aucun point n'est trouvé, retourne le nom de fichier tel quel
		return fileName
	}
	// Retourne le nom de fichier sans l'extension
	return fileName[:lastDotIndex]
}

                                                                                                                                                                                      


func AddServices(w http.ResponseWriter, r *http.Request) {
	// Structure pour la réponse
	type Response struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Data    *models.Athlete `json:"data,omitempty"`
	}

	// Récupérer le fichier image
	file, header, err := r.FormFile("image")
	if err != nil {
		logs.Errorf("Erreur lors de la récupération de l'image: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Image manquante ou invalide",
		})
		return
	}
	defer file.Close()

	// Décoder le corps de la requête pour les autres données
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

	// Téléverser l'image vers Supabase Storage
	filePath := fmt.Sprintf("athletes/%s/%s", body.Nom, header.Filename)
	uploadedURL, err := uploadImageToSupabase(file, filePath)
	if err != nil {
		logs.Errorf("Erreur lors du téléversement de l'image: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Erreur lors du téléversement de l'image",
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

	// Préparer la requête SQL pour insérer l'athlète et l'URL de l'image
	query := `INSERT INTO athletes (nom, prenom, age, image_url) VALUES ($1, $2, $3, $4) RETURNING id, nom, prenom, age, image_url`
	var insertedAthlete models.Athlete
	err = tx.QueryRowContext(ctx, query, body.Nom, body.Prenom, body.Age, uploadedURL).Scan(
		&insertedAthlete.ID,
		&insertedAthlete.Nom,
		&insertedAthlete.Prenom,
		&insertedAthlete.Age,
		&insertedAthlete.ImageURL,
	)

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
		Message: "Athlète et image ajoutés avec succès",
		Data:    &insertedAthlete,
	})
}

// Fonction pour téléverser une image vers Supabase Storage
func uploadImageToSupabase(file io.Reader, filePath string) (string, error) {
	// Créer un fichier temporaire
	tempFile, err := os.CreateTemp("", "upload-*.jpg")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())

	// Copier le contenu du fichier dans le fichier temporaire
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return "", err
	}

	// Téléverser le fichier vers Supabase Storage
	resp, err := supabaseClient.Storage.Upload("athlete-images", filePath, tempFile)
	if err != nil {
		return "", err
	}

	// Retourner l'URL publique de l'image
	return resp.PublicURL, nil
}
