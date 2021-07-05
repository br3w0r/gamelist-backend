package controller

import (
	"net/http"

	"bitbucket.org/br3w0r/gamelist-backend/entity"
	"bitbucket.org/br3w0r/gamelist-backend/service"
	"github.com/gin-gonic/gin"
)

type GameListController interface {
	PostGame(ctx *gin.Context) // Will be replaced with gRPC calls
	GetAllGames(ctx *gin.Context)

	PostGenre(ctx *gin.Context)    // Will be replaced with gRPC calls
	GetAllGenres(ctx *gin.Context) // Will be replaced with gRPC calls

	PostPlatform(ctx *gin.Context)    // Will be replaced with gRPC calls
	GetAllPlatforms(ctx *gin.Context) // Will be replaced with gRPC calls
}

type gameListController struct {
	service service.GameListService
}

func NewGameListController(service service.GameListService) GameListController {
	return &gameListController{service}
}

func (c *gameListController) PostGame(ctx *gin.Context) {
	GenericPost(ctx, &entity.GameProperties{}, c.service.SaveGame)
}

func (c *gameListController) GetAllGames(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.service.GetAllGames())
}

func (c *gameListController) PostGenre(ctx *gin.Context) {
	GenericPost(ctx, &entity.Genre{}, c.service.SaveGenre)
}

func (c *gameListController) GetAllGenres(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.service.GetAllGenres())
}

func (c *gameListController) PostPlatform(ctx *gin.Context) {
	GenericPost(ctx, &entity.Platform{}, c.service.SavePlatform)
}

func (c *gameListController) GetAllPlatforms(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.service.GetAllPlatforms())
}
