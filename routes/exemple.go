package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nedpals/supabase-go"
	_ "github.com/lib/pq"
)

var (
	db           *sql.DB
	supabaseClient *supabase.Client
)

func main() {
	// Connexion à Supabase (PostgreSQL)
	var err error
	connStr := "postgresql://<user>:<password>@<host>:<port>/<dbname>?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialisation du client Supabase
	supabaseClient = supabase.CreateClient("https://<votre-projet>.supabase.co", "<votre-cle-supabase>")

	// Configuration de l'API avec Gin
	r := gin.Default()
	r.POST("/athletes", addAthleteWithImage)
	r.Run(":8080")
}

// Struct pour les données de l'athlète
type AthleteRequest struct {
	Name string `form:"name" binding:"required"`
	Age  int    `form:"age" binding:"required"`
}

// Handler pour ajouter un athlète avec une image
func addAthleteWithImage(c *gin.Context) {
	// Récupérer les données du formulaire (nom, âge, fichier)
	var athlete AthleteRequest
	if err := c.ShouldBind(&athlete); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Récupérer le fichier image
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}
	defer file.Close()

	// Ajouter l'athlète dans la table `athletes`
	var athleteID string
	err = db.QueryRow("INSERT INTO athletes (name, age) VALUES ($1, $2) RETURNING id", athlete.Name, athlete.Age).Scan(&athleteID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add athlete"})
		return
	}

	// Téléverser l'image vers Supabase Storage
	filePath := fmt.Sprintf("athletes/%s/%s", athleteID, header.Filename)
	uploadedURL, err := uploadImageToSupabase(file, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
		return
	}

	// Ajouter les métadonnées de l'image dans la table `athlete_images`
	_, err = db.Exec("INSERT INTO athlete_images (athlete_id, file_name, file_path) VALUES ($1, $2, $3)", athleteID, header.Filename, uploadedURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image metadata"})
		return
	}

	// Réponse réussie
	c.JSON(http.StatusOK, gin.H{
		"message":   "Athlete and image added successfully",
		"athleteID": athleteID,
		"imageURL":  uploadedURL,
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