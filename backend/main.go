package main

import (
	"log"
	"net/http"
	"os"

	"costcountintimeapp/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/calculate", handlers.CORS(handlers.CalculateHandler))
	fs := http.Dir(".")
	http.HandleFunc("/", handlers.StaticHandler(fs))

	log.Printf("🚀 Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
