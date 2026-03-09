package main

import (
	"net/http"
	"costcountintimeapp/handlers"
)

func init() {
	// Initialize in-memory karma storage for Vercel
	handlers.InitKarmaDBSimple("")
}

// Handler is the main entry point for Vercel serverless functions
func Handler(w http.ResponseWriter, r *http.Request) {
	// Enable CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route requests to appropriate handlers
	switch r.URL.Path {
	case "/api/calculate":
		handlers.CalculateHandler(w, r)
	case "/api/karma/update":
		handlers.KarmaUpdateHandlerSimple(w, r)
	case "/api/karma/state":
		handlers.KarmaStateHandlerSimple(w, r)
	case "/api/karma/reset":
		handlers.KarmaResetHandlerSimple(w, r)
	case "/api/karma/delete_event":
		handlers.KarmaDeleteEventHandlerSimple(w, r)
	default:
		// Serve frontend
		http.ServeFile(w, r, "frontend/dist/index.html")
	}
}
