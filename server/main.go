package server

import (
	"log"
	"net/http"

	"github.com/br3w0r/gamelist-backend/controller"
	"github.com/br3w0r/gamelist-backend/repository"
	"github.com/br3w0r/gamelist-backend/service"
	test "github.com/br3w0r/gamelist-backend/test/stress"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/logger"
)

type ServerOptions struct {
	Production         bool
	ServeStatic        bool
	ForceScrape        bool
	StaticDir          string
	DatabaseDist       string
	ScraperGRPCAddress string
	ScraperAsync       bool
	StressTest         bool
	StressTestOptions  []string
	SilentMode         bool
	DBConfig           *repository.DBConfig
}

func NewServer(options ServerOptions) *gin.Engine {
	var (
		// DB dialector init
		dialector = repository.NewDBDialector(options.DBConfig)

		// Repos
		gamelistRepository repository.GamelistRepository = repository.NewGamelistRepository(
			options.DatabaseDist, dialector,
			logger.Config{
				Colorful: true,
				IgnoreRecordNotFoundError: true,
				LogLevel: logger.Error,
			},
		)

		// Services
		gamelistService service.GameListService = service.NewGameListService(gamelistRepository, options.ScraperGRPCAddress)
		jwtService      service.JWTService      = service.NewJWTService(gamelistRepository)

		// Controllers
		gamelistController controller.GameListController = controller.NewGameListController(gamelistService, jwtService)
	)

	if options.ForceScrape {
		log.Println("Force scraping.")

		if options.ScraperAsync {
			go gamelistService.ScrapeGames()
		} else {
			gamelistService.ScrapeGames()
		}
	}

	// First version. Should be remade
	if options.StressTest {
		test.RunStress(gamelistRepository, options.StressTestOptions)
		return nil
	}

	if options.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	server := gin.New()

	server.Use(gin.Recovery())
	if !options.SilentMode {
		server.Use(gin.Logger())
	}

	if options.ServeStatic {
		server.Static("/css", options.StaticDir+"/css")
		server.Static("/js", options.StaticDir+"/js")
		server.LoadHTMLGlob(options.StaticDir + "/*.html")

		// For SPA on vue
		server.NoRoute(func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "index.html", nil)
		})
	}

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
		if options.Production {
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

	return server
}
