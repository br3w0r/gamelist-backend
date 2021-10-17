package controller

import (
	"net/http"
	"strings"

	"github.com/br3w0r/gamelist-backend/entity"
	"github.com/br3w0r/gamelist-backend/service"
	utilErrs "github.com/br3w0r/gamelist-backend/util/errors"
	"github.com/gin-gonic/gin"
)

type GameListController interface {
	GetAllGames(ctx *gin.Context)
	GetAllGamesTyped(ctx *gin.Context)
	GetMyGameList(ctx *gin.Context)
	SearchGames(ctx *gin.Context)
	GameDetails(ctx *gin.Context)

	AcquireJWTPair(ctx *gin.Context)
	RefreshJWTPair(ctx *gin.Context)
	RevokeRefreshToken(ctx *gin.Context)
	DeleteAllRefreshTokens(ctx *gin.Context)
	Authorized(ctx *gin.Context)

	// Will be replaced with gRPC calls
	PostGame(ctx *gin.Context)

	PostListType(ctx *gin.Context)
	GetAllListTypes(ctx *gin.Context)
	ListGame(ctx *gin.Context)

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
	games, err := c.gamelistService.GetAllGames()
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, games)
}

func (c *gameListController) GetAllGamesTyped(ctx *gin.Context) {
	nickname := ctx.MustGet("nickname").(string)
	var request entity.GameBatchRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ErrorSender(ctx, utilErrs.JSONParseErr(err))
		return
	}

	games, err := c.gamelistService.GetAllGamesTyped(nickname, request.Last, request.BatchSize)
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, games)
}

func (c *gameListController) GetMyGameList(ctx *gin.Context) {
	nickname := ctx.MustGet("nickname").(string)

	games, err := c.gamelistService.GetUserGameList(nickname)
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, games)
}

func (c *gameListController) SearchGames(ctx *gin.Context) {
	var request entity.SearchRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ErrorSender(ctx, utilErrs.JSONParseErr(err))
		return
	}

	games, err := c.gamelistService.SearchGames(request.Name)
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, games)
}

func (c *gameListController) GameDetails(ctx *gin.Context) {
	nickname := ctx.MustGet("nickname").(string)
	var request entity.GameDetailsRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ErrorSender(ctx, utilErrs.JSONParseErr(err))
		return
	}

	gameDetails, err := c.gamelistService.GetGameDetails(nickname, request.Id)
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gameDetails)
}

func (c *gameListController) PostListType(ctx *gin.Context) {
	GenericPost(ctx, &entity.ListType{}, c.gamelistService.CreateListType)
}

func (c *gameListController) GetAllListTypes(ctx *gin.Context) {
	types, err := c.gamelistService.GetAllListTypes()
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, types)
}

func (c *gameListController) ListGame(ctx *gin.Context) {
	nickname := ctx.MustGet("nickname").(string)

	var gameList entity.GameListRequest
	err := ctx.ShouldBindJSON(&gameList)
	if err != nil {
		ErrorSender(ctx, utilErrs.JSONParseErr(err))
		return
	}

	err = c.gamelistService.ListGame(nickname, gameList.GameId, gameList.ListType)
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ResponseOK(ctx)
}

func (c *gameListController) PostGenre(ctx *gin.Context) {
	GenericPost(ctx, &entity.Genre{}, c.gamelistService.SaveGenre)
}

func (c *gameListController) GetAllGenres(ctx *gin.Context) {
	genres, err := c.gamelistService.GetAllGenres()
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, genres)
}

func (c *gameListController) PostPlatform(ctx *gin.Context) {
	GenericPost(ctx, &entity.Platform{}, c.gamelistService.SavePlatform)
}

func (c *gameListController) GetAllPlatforms(ctx *gin.Context) {
	platforms, err := c.gamelistService.GetAllPlatforms()
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, platforms)
}

func (c *gameListController) PostProfile(ctx *gin.Context) {
	GenericPost(ctx, &entity.Profile{}, c.gamelistService.CreateProfile)
}

func (c *gameListController) GetAllProfiles(ctx *gin.Context) {
	profiles, err := c.gamelistService.GetAllProfiles()
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, profiles)
}

func (c *gameListController) AcquireJWTPair(ctx *gin.Context) {
	var login entity.LoginProfile
	err := ctx.ShouldBindJSON(&login)
	if err != nil {
		ErrorSender(ctx, utilErrs.JSONParseErr(err))
		return
	}

	profile, err := c.gamelistService.CheckLogin(login)
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	pair, err := c.jwtService.GenerateTokens(profile.Nickname)
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, pair)
}

func (c *gameListController) RefreshJWTPair(ctx *gin.Context) {
	var refresh entity.RefreshRequest
	err := ctx.ShouldBindJSON(&refresh)
	if err != nil {
		ErrorSender(ctx, utilErrs.JSONParseErr(err))
		return
	}

	pair, err := c.jwtService.RefreshTokens(refresh.RefreshToken)
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, pair)
}

func (c *gameListController) RevokeRefreshToken(ctx *gin.Context) {
	var request entity.RefreshRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ErrorSender(ctx, utilErrs.JSONParseErr(err))
		return
	}

	err = c.jwtService.RevokeRefreshToken(request.RefreshToken)
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ResponseOK(ctx)
}

func (c *gameListController) DeleteAllRefreshTokens(ctx *gin.Context) {
	nickname := ctx.MustGet("nickname").(string)
	err := c.jwtService.DeleteAllUserRefreshTokens(nickname)
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ResponseOK(ctx)
}

func (c *gameListController) Authorized(ctx *gin.Context) {
	var token string
	authHeader, ok := ctx.Request.Header["Authorization"]
	if !ok {
		ErrorSender(ctx, utilErrs.New(utilErrs.Unauthorized, nil, "no authorization header provided"))
		return
	}

	list := strings.Split(authHeader[0], " ")
	if len(list) != 2 || list[0] != "Bearer" {
		ErrorSender(ctx, utilErrs.New(utilErrs.Unauthorized, nil, "wrong authorization header format"))
		return
	}

	token = list[1]
	nickname, err := c.jwtService.Authenticate(token)
	if err != nil {
		ErrorSender(ctx, utilErrs.New(utilErrs.Unauthorized, nil, "authentication failed"))
		return
	}

	ctx.Set("nickname", nickname)
}

func (c *gameListController) PostSocialType(ctx *gin.Context) {
	GenericPost(ctx, &entity.SocialType{}, c.gamelistService.SaveSocialType)
}

func (c *gameListController) GetAllSocialtypes(ctx *gin.Context) {
	types, err := c.gamelistService.GetAllSocialTypes()
	if err != nil {
		ErrorSender(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, types)
}
