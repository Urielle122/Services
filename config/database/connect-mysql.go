package database

import (
	"database/sql"
	"fmt"
	"services/logs"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

func ConnectToSupabase() (*sql.DB, error) {
	// Paramètres de connexion
	const (
		host     = "aws-0-ca-central-1.pooler.supabase.com"
		port     = 6543
		user     = "postgres.yokfzzrreskhcvpclcib"
		password = "xO5xAWXR0mTlRUBp" // Remplacez par votre mot de passe
		dbname   = "postgres"
	)

	// Construction de la chaîne de connexion PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	// Ouverture de la connexion
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logs.Errorf("Erreur lors de l'ouverture de la connexion à Supabase : %v", err)
		return nil, err
	}

	// Test de la connexion
	err = db.Ping()
	if err != nil {
		logs.Errorf("Erreur lors du test de connexion à Supabase : %v", err)
		return nil, err
	}

	logs.Info("Connecté avec succès à Supabase")
	return db, nil
}

// Fonction optionnelle pour fermer la connexion
func CloseConnection(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			logs.Errorf("Erreur lors de la fermeture de la connexion : %v", err)
		}
	}
}