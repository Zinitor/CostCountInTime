package handler

import (
	"log"
	"net/http"
	"os"
	"sync"

	"costcountintimeapp/application/service"
	"costcountintimeapp/domain/repository"
	"costcountintimeapp/infrastructure/persistence"
	interfaceshttp "costcountintimeapp/interfaces/http"
	"costcountintimeapp/interfaces/router"
)

var (
	routerInstance *router.Router
	routerOnce     sync.Once
)

// init initializes storage for Vercel serverless functions
func init() {
	routerOnce.Do(func() {
		// Try to use Vercel Blob in production, fallback to in-memory
		var karmaRepo repository.KarmaRepository
		var err error

		if os.Getenv("BLOB_READ_WRITE_TOKEN") != "" {
			karmaRepo, err = persistence.NewBlobKarmaRepository()
			if err != nil {
				log.Printf("⚠️ Blob storage failed: %v, falling back to in-memory", err)
				karmaRepo = persistence.NewInMemoryKarmaRepository()
			} else {
				log.Println("✅ Using Vercel Blob for karma storage")
			}
		} else {
			log.Println("⚠️ BLOB_READ_WRITE_TOKEN not set, using in-memory karma storage")
			karmaRepo = persistence.NewInMemoryKarmaRepository()
		}

		tagRepo := persistence.NewInMemoryTagRepository()

		calculateService := service.NewCalculateService(tagRepo)
		karmaService := service.NewKarmaService(karmaRepo, tagRepo)

		calculateHandler := interfaceshttp.NewCalculateHandler(calculateService)
		karmaHandler := interfaceshttp.NewKarmaHandler(karmaService)
		tagHandler := interfaceshttp.NewTagHandler(tagRepo)
		staticHandler := interfaceshttp.StaticHandler(http.Dir("."))

		routerInstance = router.NewRouter(calculateHandler, karmaHandler, tagHandler, staticHandler)
	})
}

// Handler is the main entry point for Vercel serverless functions
func Handler(w http.ResponseWriter, r *http.Request) {
	routerInstance.ServeHTTP(w, r)
}
