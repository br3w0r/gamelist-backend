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
	var game entity.GameProperties
	err := ctx.ShouldBindJSON(&game)

	if err != nil {
		ErrorSender(ctx, err)
	} else {
		err = c.service.SaveGame(game)
		if err != nil {
			ErrorSender(ctx, err)
		} else {
			ResponseOK(ctx)
		}
	}
}

func (c *gameListController) GetAllGames(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.service.GetAllGames())
}

func (c *gameListController) PostGenre(ctx *gin.Context) {
	var genre entity.Genre
	err := ctx.ShouldBindJSON(&genre)

	if err != nil {
		ErrorSender(ctx, err)
	} else {
		err = c.service.SaveGenre(genre)
		if err != nil {
			ErrorSender(ctx, err)
		} else {
			ResponseOK(ctx)
		}
	}
}

func (c *gameListController) GetAllGenres(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.service.GetAllGenres())
}

func (c *gameListController) PostPlatform(ctx *gin.Context) {
	var platform entity.Platform
	err := ctx.ShouldBindJSON(&platform)

	if err != nil {
		ErrorSender(ctx, err)
	} else {
		err = c.service.SavePlatform(platform)
		if err != nil {
			ErrorSender(ctx, err)
		} else {
			ResponseOK(ctx)
		}
	}
}

func (c *gameListController) GetAllPlatforms(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.service.GetAllPlatforms())
}
