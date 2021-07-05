package main

import (
	"flag"
	"log"
	"net/http"

	"bitbucket.org/br3w0r/gamelist-backend/controller"
	"bitbucket.org/br3w0r/gamelist-backend/repository"
	"bitbucket.org/br3w0r/gamelist-backend/service"
	"github.com/gin-gonic/gin"
)

func main() {
	var (
		migrate     *bool = flag.Bool("migrate", false, "Force AutoMigrate if true")
		forceScrape *bool = flag.Bool("force-scrape", false, "Force all games scraping")
	)
	flag.Parse()

	if *migrate {
		log.Println("Force migration.")
	}

	var (
		gamelistRepository                               = repository.NewGamelistRepository("gamelist.db", *migrate)
		gamelistService    service.GameListService       = service.NewGameListService(gamelistRepository)
		gamelistController controller.GameListController = controller.NewGameListController(gamelistService)
	)

	if *forceScrape {
		log.Println("Force scraping.")
		go gamelistService.ScrapeGames()
	}

	server := gin.Default()

	server.Static("/css", "../gamelist-frontend/gamelist/dist/css")
	server.Static("/js", "../gamelist-frontend/gamelist/dist/js")
	server.LoadHTMLGlob("../gamelist-frontend/gamelist/dist/*.html")

	// For SPA on vue
	server.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	apiRoutes := server.Group("/api/v0")
	{
		apiRoutes.GET("/games/all", gamelistController.GetAllGames)

		// This will be replaced with gRPC call
		apiRoutes.POST("/games", gamelistController.PostGame)

		apiRoutes.GET("/genres", gamelistController.GetAllGenres)
		apiRoutes.POST("/genres", gamelistController.PostGenre)

		apiRoutes.GET("/platforms", gamelistController.GetAllPlatforms)
		apiRoutes.POST("/platforms", gamelistController.PostPlatform)

		apiRoutes.GET("/profiles", gamelistController.GetAllProfiles)
		apiRoutes.POST("/profiles", gamelistController.PostProfile)

		apiRoutes.POST("/login-check", gamelistController.CheckLogin)

		apiRoutes.GET("/social-types", gamelistController.GetAllSocialtypes)
		apiRoutes.POST("/social-types", gamelistController.PostSocialType)
	}

	server.Run(":8080")
}
