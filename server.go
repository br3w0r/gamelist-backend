package main

import (
	"log"
	"net/http"

	"bitbucket.org/br3w0r/gamelist-backend/controller"
	"bitbucket.org/br3w0r/gamelist-backend/helpers"
	"bitbucket.org/br3w0r/gamelist-backend/repository"
	"bitbucket.org/br3w0r/gamelist-backend/service"
	"github.com/gin-gonic/gin"
)

var (
	PRODUCTION_MODE      string = helpers.GetEnvOrDefault("PRODUCTION_MODE", "0")
	SERVE_STATIC         string = helpers.GetEnvOrDefault("SERVE_STATIC", "1")
	FORCE_MIGRATE        string = helpers.GetEnvOrDefault("FORCE_MIGRATE", "0")
	FORCE_SCRAPE         string = helpers.GetEnvOrDefault("FORCE_SCRAPE", "0")
	STATIC_DIR           string = helpers.GetEnvOrDefault("STATIC_FOLDER", "../gamelist-frontend/gamelist/dist")
	DATABASE_DIR         string = helpers.GetEnvOrDefault("DATABASE_DIR", ".")
	SCRAPER_GRPC_ADDRESS string = helpers.GetEnvOrDefault("SCRAPER_GRPC_ADDRESS", "localhost")
)

func main() {
	if FORCE_MIGRATE == "1" {
		log.Println("Force migration.")
	}

	var (
		// Repos
		gamelistRepository = repository.NewGamelistRepository(DATABASE_DIR+"/gamelist.db", FORCE_MIGRATE == "1")

		// Services
		gamelistService service.GameListService = service.NewGameListService(gamelistRepository, SCRAPER_GRPC_ADDRESS)
		jwtService      service.JWTService      = service.NewJWTService(gamelistRepository)

		// Controllers
		gamelistController controller.GameListController = controller.NewGameListController(gamelistService, jwtService)
	)

	if FORCE_SCRAPE == "1" {
		log.Println("Force scraping.")
		go gamelistService.ScrapeGames()
	}

	if PRODUCTION_MODE == "1" {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.Default()

	if SERVE_STATIC == "1" {
		server.Static("/css", STATIC_DIR+"/css")
		server.Static("/js", STATIC_DIR+"/js")
		server.LoadHTMLGlob(STATIC_DIR + "/*.html")
	}

	// For SPA on vue
	server.NoRoute(func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	apiRoutes := server.Group("/api/v0")
	{
		apiRoutes.POST("/games/all",
			gamelistController.Authorized,
			gamelistController.GetAllGamesTyped,
		)

		apiRoutes.POST("/list-game",
			gamelistController.Authorized,
			gamelistController.ListGame,
		)

		apiRoutes.GET("/my-games",
			gamelistController.Authorized,
			gamelistController.GetMyGameList,
		)

		apiRoutes.POST("/games/search",
			gamelistController.Authorized,
			gamelistController.SearchGames,
		)

		apiRoutes.POST("/games/details",
			gamelistController.Authorized,
			gamelistController.GameDetails,
		)

		apiRoutes.POST("/profiles", gamelistController.PostProfile)

		apiRoutes.POST("/aquire-tokens", gamelistController.AcquireJWTPair)
		apiRoutes.POST("/refresh-tokens", gamelistController.RefreshJWTPair)
		apiRoutes.POST("/revoke-token", gamelistController.RevokeRefreshToken)
		apiRoutes.GET("/delete-all-refresh-tokens",
			gamelistController.Authorized,
			gamelistController.DeleteAllRefreshTokens,
		)

		// This will be replaced with gRPC admin shell
		if PRODUCTION_MODE == "0" {
			apiRoutes.POST("/games", gamelistController.PostGame)

			apiRoutes.POST("/list-types", gamelistController.PostListType)
			apiRoutes.GET("/list-types", gamelistController.GetAllListTypes)

			apiRoutes.GET("/genres", gamelistController.GetAllGenres)
			apiRoutes.POST("/genres", gamelistController.PostGenre)

			apiRoutes.GET("/platforms", gamelistController.GetAllPlatforms)
			apiRoutes.POST("/platforms", gamelistController.PostPlatform)

			apiRoutes.GET("/profiles", gamelistController.GetAllProfiles)

			apiRoutes.GET("/social-types", gamelistController.GetAllSocialtypes)
			apiRoutes.POST("/social-types", gamelistController.PostSocialType)
		}
	}

	server.Run(":8080")
}
