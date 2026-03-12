package handler

import (
	"net/http"
	"sync"

	"costcountintimeapp/application/service"
	"costcountintimeapp/infrastructure/persistence"
	interfaceshttp "costcountintimeapp/interfaces/http"
	"costcountintimeapp/interfaces/router"
)

var (
	routerInstance *router.Router
	routerOnce     sync.Once
)

// init initializes in-memory karma storage for Vercel
func init() {
	routerOnce.Do(func() {
		// Для serverless используем in-memory репозиторий
		karmaRepo := persistence.NewInMemoryKarmaRepository()
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
