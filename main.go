package main

import (
	"net/http"
	"log"

	"services/core"
	"services/routes"

)

func main() {
	core.InitConnection()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /add", routes.AddServices)
	mux.HandleFunc("GET /service", routes.GetServices)
	log.Println("Serveur démarré sur le port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}