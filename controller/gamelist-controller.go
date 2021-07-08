package controller

import (
	"net/http"
	"strings"

	"bitbucket.org/br3w0r/gamelist-backend/entity"
	"bitbucket.org/br3w0r/gamelist-backend/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	AcquireJWTPair(ctx *gin.Context)
	RefreshJWTPair(ctx *gin.Context)
	DeleteAllRefreshTokens(ctx *gin.Context)
	Authorized(ctx *gin.Context)

	PostSocialType(ctx *gin.Context)
	GetAllSocialtypes(ctx *gin.Context)
}

type gameListController struct {
	gamelistService service.GameListService
	jwtService      service.JWTService
}

func NewGameListController(gamelistService service.GameListService, jwtService service.JWTService) GameListController {
	return &gameListController{
		gamelistService: gamelistService,
		jwtService:      jwtService,
	}
}

func (c *gameListController) PostGame(ctx *gin.Context) {
	GenericPost(ctx, &entity.GameProperties{}, c.gamelistService.SaveGame)
}

func (c *gameListController) GetAllGames(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.gamelistService.GetAllGames())
}

func (c *gameListController) PostGenre(ctx *gin.Context) {
	GenericPost(ctx, &entity.Genre{}, c.gamelistService.SaveGenre)
}

func (c *gameListController) GetAllGenres(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.gamelistService.GetAllGenres())
}

func (c *gameListController) PostPlatform(ctx *gin.Context) {
	GenericPost(ctx, &entity.Platform{}, c.gamelistService.SavePlatform)
}

func (c *gameListController) GetAllPlatforms(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.gamelistService.GetAllPlatforms())
}

func (c *gameListController) PostProfile(ctx *gin.Context) {
	GenericPost(ctx, &entity.Profile{}, c.gamelistService.CreateProfile)
}

func (c *gameListController) GetAllProfiles(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.gamelistService.GetAllProfiles())
}

func (c *gameListController) AcquireJWTPair(ctx *gin.Context) {
	var login entity.LoginProfile
	err := ctx.ShouldBindJSON(&login)
	if err != nil {
		ErrorSender(ctx, err)
	} else {
		profile, err := c.gamelistService.CheckLogin(login)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
		} else {
			pair, err := c.jwtService.GenerateTokens(profile.Nickname)
			if err != nil {
				ErrorSender(ctx, err)
			} else {
				ctx.JSON(http.StatusOK, pair)
			}
		}
	}
}

func (c *gameListController) RefreshJWTPair(ctx *gin.Context) {
	var refresh entity.RefreshRequest
	err := ctx.ShouldBindJSON(&refresh)
	if err != nil {
		ErrorSender(ctx, err)
	} else {
		pair, err := c.jwtService.RefreshTokens(refresh.RefreshToken)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			} else {
				ErrorSender(ctx, err)
			}
		} else {
			ctx.JSON(http.StatusOK, pair)
		}
	}
}

func (c *gameListController) DeleteAllRefreshTokens(ctx *gin.Context) {
	nickname := ctx.MustGet("nickname").(string)
	err := c.jwtService.DeleteAllUserRefreshTokens(nickname)
	if err != nil {
		ErrorSender(ctx, err)
	} else {
		ResponseOK(ctx)
	}
}

func (c *gameListController) Authorized(ctx *gin.Context) {
	token := strings.Split(ctx.Request.Header["Authorization"][0], " ")[1]
	nickname, err := c.jwtService.Authenticate(token)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
	} else {
		ctx.Set("nickname", nickname)
	}
}

func (c *gameListController) PostSocialType(ctx *gin.Context) {
	GenericPost(ctx, &entity.SocialType{}, c.gamelistService.SaveSocialType)
}

func (c *gameListController) GetAllSocialtypes(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.gamelistService.GetAllSocialTypes())
}
