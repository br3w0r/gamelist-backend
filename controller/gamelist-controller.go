package controller

import (
	"net/http"

	"bitbucket.org/br3w0r/gamelist-backend/entity"
	"bitbucket.org/br3w0r/gamelist-backend/service"
	"github.com/gin-gonic/gin"
)

type GameListController interface {
	GetAllGames(ctx *gin.Context)

	// Will be replaced with gRPC calls
	PostGame(ctx *gin.Context)

	PostGenre(ctx *gin.Context)
	GetAllGenres(ctx *gin.Context)

	PostPlatform(ctx *gin.Context)
	GetAllPlatforms(ctx *gin.Context)

	PostProfile(ctx *gin.Context)
	GetAllProfiles(ctx *gin.Context)

	PostSocialType(ctx *gin.Context)
	GetAllSocialtypes(ctx *gin.Context)
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

func (c *gameListController) PostProfile(ctx *gin.Context) {
	GenericPost(ctx, &entity.ProfileInfo{}, c.service.SaveProfile)
}

func (c *gameListController) GetAllProfiles(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.service.GetAllProfiles())
}

func (c *gameListController) PostSocialType(ctx *gin.Context) {
	GenericPost(ctx, &entity.SocialType{}, c.service.SaveSocialType)
}

func (c *gameListController) GetAllSocialtypes(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.service.GetAllSocialTypes())
}
