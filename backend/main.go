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

	// Инициализация Karma DB
	if err := handlers.InitKarmaDB("karma.db"); err != nil {
		log.Printf("⚠️ Failed to init karma DB: %v", err)
	} else {
		log.Println("✅ Karma DB initialized")
	}

	// API endpoints
	http.HandleFunc("/api/calculate", handlers.CORS(handlers.CalculateHandler))
	http.HandleFunc("/api/karma/update", handlers.CORS(handlers.KarmaUpdateHandler))
	http.HandleFunc("/api/karma/state", handlers.CORS(handlers.KarmaStateHandler))
	http.HandleFunc("/api/karma/reset", handlers.CORS(handlers.KarmaResetHandler))
	http.HandleFunc("/api/karma/delete_event", handlers.CORS(handlers.KarmaDeleteEventHandler))

	// Статические файлы (frontend)
	fs := http.Dir(".")
	http.HandleFunc("/", handlers.StaticHandler(fs))

	log.Printf("🚀 Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
