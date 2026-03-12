package handler

import (
	"log"
	"net/http"
	"os"

	"costcountintimeapp/application/service"
	"costcountintimeapp/infrastructure/persistence"
	interfaceshttp "costcountintimeapp/interfaces/http"
	"costcountintimeapp/interfaces/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Инициализация инфраструктуры (Infrastructure Layer)
	karmaRepo, err := persistence.NewSQLiteKarmaRepository("karma.db")
	if err != nil {
		log.Printf("⚠️ Failed to init karma DB: %v", err)
	} else {
		log.Println("✅ Karma DB initialized")
	}
	defer karmaRepo.Close()

	tagRepo := persistence.NewInMemoryTagRepository()
	log.Println("✅ Tag repository initialized")

	// Инициализация сервисов приложения (Application Layer)
	calculateService := service.NewCalculateService(tagRepo)
	karmaService := service.NewKarmaService(karmaRepo, tagRepo)

	// Инициализация HTTP обработчиков (Interfaces Layer)
	calculateHandler := interfaceshttp.NewCalculateHandler(calculateService)
	karmaHandler := interfaceshttp.NewKarmaHandler(karmaService)
	tagHandler := interfaceshttp.NewTagHandler(tagRepo)
	staticHandler := interfaceshttp.StaticHandler(http.Dir("."))

	// Инициализация роутера
	r := router.NewRouter(calculateHandler, karmaHandler, tagHandler, staticHandler)

	log.Printf("🚀 Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
