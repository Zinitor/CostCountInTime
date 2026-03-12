package http

import (
	"net/http"
	"path/filepath"
)

// StaticHandler создаёт обработчик для раздачи статических файлов
func StaticHandler(fs http.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		filePath := filepath.Clean("frontend/dist" + path)
		f, err := fs.Open(filePath)
		if err != nil {
			http.ServeFile(w, r, "frontend/dist/index.html")
			return
		}
		defer f.Close()
		http.ServeFile(w, r, filePath)
	}
}
