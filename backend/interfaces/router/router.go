package router

import (
	"net/http"

	httpint "costcountintimeapp/interfaces/http"
)

// Router управляет маршрутизацией HTTP запросов
type Router struct {
	mux *http.ServeMux
}

// NewRouter создаёт новый Router
func NewRouter(
	calculateHandler *httpint.CalculateHandler,
	karmaHandler *httpint.KarmaHandler,
	tagHandler *httpint.TagHandler,
	staticHandler http.HandlerFunc,
) *Router {
	mux := http.NewServeMux()

	// API endpoints с CORS middleware
	mux.HandleFunc("/api/calculate", httpint.CORS(calculateHandler.ServeHTTP))
	mux.HandleFunc("/api/karma/update", httpint.CORS(karmaHandler.UpdateHandler))
	mux.HandleFunc("/api/karma/state", httpint.CORS(karmaHandler.StateHandler))
	mux.HandleFunc("/api/karma/reset", httpint.CORS(karmaHandler.ResetHandler))
	mux.HandleFunc("/api/karma/delete_event", httpint.CORS(karmaHandler.DeleteEventHandler))
	mux.HandleFunc("/api/tags", httpint.CORS(tagHandler.ServeHTTP))

	// Статические файлы (frontend)
	mux.HandleFunc("/", staticHandler)

	return &Router{
		mux: mux,
	}
}

// ServeHTTP реализует http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
