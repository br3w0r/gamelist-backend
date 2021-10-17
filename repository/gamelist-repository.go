package repository

import (
	"fmt"
	"log"

	"github.com/br3w0r/gamelist-backend/entity"
	utilErrs "github.com/br3w0r/gamelist-backend/util/errors"
	utilLogger "github.com/br3w0r/gamelist-backend/util/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type GamelistRepository interface {
	SaveGame(game entity.GameProperties) error
	GetAllGames() ([]entity.GameProperties, error)
	GetAllGamesTyped(nickname string, last uint64, batchSize int) ([]entity.TypedGameListProperties, error)
	GetUserGameList(nickname string) ([]entity.TypedGameListProperties, error)
	SearchGames(name string) ([]entity.GameSearchResult, error)
	GetGameDetails(nickname string, id uint64) (*entity.GameDetailsResponse, error)

	CreateListType(listType entity.ListType) error
	GetAllListTypes() ([]entity.ListType, error)
	ListGame(nickname string, gameId uint64, listType uint64) error

	SaveGenre(genre entity.Genre) error
	GetAllGenres() ([]entity.Genre, error)

	SavePlatform(platform entity.Platform) error
	GetAllPlatforms() ([]entity.Platform, error)

	CreateProfile(profile entity.Profile) error
	SaveProfile(profile entity.Profile) error
	GetAllProfiles() ([]entity.ProfileInfo, error)
	GetProfile(login entity.ProfileCreds) (*entity.Profile, error)

	SaveRefreshToken(nickname string, tokenString string) error
	FindRefreshToken(nickname string, tokenString string) error
	DeleteRefreshToken(tokenString string) error
	DeleteAllUserRefreshTokens(nickname string) error

	SaveSocialType(socialType entity.SocialType) error
	GetAllSocialTypes() ([]entity.SocialType, error)
}

type gameListRepository struct {
	db *gorm.DB
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
	SSL      bool
	TimeZone string
}

var (
	GAMES_BATCH_SIZE_LIMIT int = 10
	ErrDbConnection            = "Failed to connect database."
)

func NewDBDialector(conf *DBConfig) gorm.Dialector {
	var sslString string
	if conf.SSL {
		sslString = "require"
	} else {
		sslString = "disable"
	}

	dsn := fmt.Sprint("host=", conf.Host,
		" user=", conf.User,
		" password=", conf.Password,
		" dbname=", conf.DBName,
		" port=", conf.Port,
		" sslmode=", sslString,
		" TimeZone=", conf.TimeZone,
	)

	return postgres.Open(dsn)
}

func NewGamelistRepository(dbName string, dialector gorm.Dialector, loggerConf logger.Config) GamelistRepository {
	var (
		db  *gorm.DB
		err error
	)

	logger := logger.New(
		log.New(utilLogger.Logger, "", log.LstdFlags),
		loggerConf,
	)

	db, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		panic(ErrDbConnection)
	}

	return &gameListRepository{
		db: db,
	}
}

func (r *gameListRepository) SaveGame(game entity.GameProperties) error {
	for i := range game.Platforms {
		res := r.db.First(&game.Platforms[i], game.Platforms[i])
		if res.Error != nil {
			return utilErrs.FromGORM(res,
				fmt.Sprintf("couldn't find platform with id %d and name %s",
					game.Platforms[i].ID, game.Platforms[i].Name,
				))
		}
	}
	for i := range game.Genres {
		res := r.db.First(&game.Genres[i], game.Genres[i])
		if res.Error != nil {
			return utilErrs.FromGORM(res,
				fmt.Sprintf("couldn't find genre with id %d and name %s",
					game.Genres[i].ID, game.Genres[i].Name,
				))
		}
	}

	res := r.db.Save(&game)
	if res.Error != nil {
		return utilErrs.FromGORM(res, "failed to save game")
	}

	return nil
}

func (r *gameListRepository) GetAllGames() ([]entity.GameProperties, error) {
	var games []entity.GameProperties
	res := r.db.Preload(clause.Associations).Find(&games)
	if res.Error != nil {
		return nil, utilErrs.FromGORM(res, "failed to get games")
	}

	return games, nil
}

func (r *gameListRepository) GetAllGamesTyped(nickname string, last uint64, batchSize int) ([]entity.TypedGameListProperties, error) {
	if batchSize > GAMES_BATCH_SIZE_LIMIT {
		batchSize = GAMES_BATCH_SIZE_LIMIT
	}

	userId, err := r.findUserIDByNickname(nickname)
	if err != nil {
		return nil, err
	}

	var games []entity.TypedGameListProperties
	res := r.db.Table("game_properties").
		Joins("left join profile_game on game_properties.id = profile_game.game_id and profile_game.profile_id = ?", userId).
		Where("game_properties.id > ?", last).
		Limit(batchSize).
		Scan(&games)

	if res.Error != nil {
		return nil, utilErrs.FromGORM(res, "failed to get games")
	}

	return games, nil
}

func (r *gameListRepository) GetUserGameList(nickname string) ([]entity.TypedGameListProperties, error) {
	var games []entity.TypedGameListProperties
	res := r.db.Table("game_properties").Select(
		"game_properties.id, game_properties.name, game_properties.image_url, game_properties.year_released, profile_game.list_type_id",
	).Joins(
		"join profile_game on game_properties.id = profile_game.game_id and profile_game.list_type_id != 0",
	).Joins(
		"join profile on profile_game.profile_id = profile.id and profile.nickname = ?",
		nickname,
	).Scan(&games)

	if res.Error != nil {
		return nil, utilErrs.FromGORM(res, "failed to get user game list")
	}

	return games, nil
}

func (r *gameListRepository) SearchGames(name string) ([]entity.GameSearchResult, error) {
	var games []entity.GameSearchResult
	if len(name) > 1 {
		res := r.db.Table("game_properties").
			Where("name LIKE ?", name+"%").
			Limit(10).
			Find(&games)

		if res.Error != nil {
			return nil, utilErrs.FromGORM(res, "failed to search games")
		}
	}
	return games, nil
}

func (r *gameListRepository) GetGameDetails(nickname string, gameId uint64) (*entity.GameDetailsResponse, error) {
	userId, err := r.findUserIDByNickname(nickname)
	if err != nil {
		return nil, err
	}

	var gameDetails entity.GameDetailsResponse
	res := r.db.Table("game_properties").
		Joins("left join profile_game on game_properties.id = profile_game.game_id and profile_game.profile_id = ?", userId).
		Where("game_properties.id = ?", gameId).
		Limit(1).
		Scan(&(gameDetails.Game))

	if res.Error != nil || res.RowsAffected == 0 {
		return nil, utilErrs.FromGORM(res, "failed to get game")
	}

	res = r.db.Table("platform").Select("platform.name").
		Joins("inner join game_platforms on game_platforms.game_properties_id = ? and game_platforms.platform_id = platform.id", gameId).
		Scan(&(gameDetails.Platforms))

	if res.Error != nil {
		return nil, utilErrs.FromGORM(res, "failed to get game's platforms")
	}

	res = r.db.Table("genre").Select("genre.name").
		Joins("inner join game_genres on game_genres.game_properties_id = ? and game_genres.genre_id = genre.id", gameId).
		Scan(&(gameDetails.Genres))

	if res.Error != nil {
		return nil, utilErrs.FromGORM(res, "failed to get game's genres")
	}

	return &gameDetails, nil
}

func (r *gameListRepository) CreateListType(listType entity.ListType) error {
	res := r.db.Create(&listType)
	if res.Error != nil {
		return utilErrs.FromGORM(res, "failed to create list type")
	}

	return nil
}

func (r *gameListRepository) GetAllListTypes() ([]entity.ListType, error) {
	var types []entity.ListType
	res := r.db.Find(&types)
	if res.Error != nil {
		return nil, utilErrs.FromGORM(res, "failed to get list types")
	}

	return types, nil
}

func (r *gameListRepository) ListGame(nickname string, gameId uint64, listType uint64) error {
	userId, err := r.findUserIDByNickname(nickname)
	if err != nil {
		return err
	}

	res := r.db.First(&entity.GameProperties{}, gameId)
	if res.Error != nil {
		return utilErrs.FromGORM(res, fmt.Sprint("couldn't find game with id: ", gameId))
	}

	if listType != 0 {
		res = r.db.First(&entity.ListType{}, listType)
		if res.Error != nil {
			return utilErrs.FromGORM(res, fmt.Sprint("couldn't find list type with id: ", listType))
		}
	}

	listGame := entity.ProfileGame{
		ProfileID:  userId,
		GameID:     gameId,
		ListTypeID: listType,
	}

	res = r.db.Save(&listGame)
	if res.Error != nil {
		return utilErrs.FromGORM(res, "failed to save changes")
	}

	return nil
}

func (r *gameListRepository) SaveGenre(genre entity.Genre) error {
	res := r.db.Save(&genre)
	if res.Error != nil {
		return utilErrs.FromGORM(res, "failed to save genre")
	}

	return nil
}

func (r *gameListRepository) GetAllGenres() ([]entity.Genre, error) {
	var genres []entity.Genre
	res := r.db.Find(&genres)
	if res.Error != nil {
		return nil, utilErrs.FromGORM(res, "failed to get genres")
	}

	return genres, nil
}

func (r *gameListRepository) SavePlatform(platform entity.Platform) error {
	res := r.db.Save(&platform)
	if res.Error != nil {
		return utilErrs.FromGORM(res, "failed to save platform")
	}

	return nil
}

func (r *gameListRepository) GetAllPlatforms() ([]entity.Platform, error) {
	var platforms []entity.Platform
	res := r.db.Find(&platforms)
	if res.Error != nil {
		return nil, utilErrs.FromGORM(res, "failed to get platforms")
	}

	return platforms, nil
}

func (r *gameListRepository) CreateProfile(profile entity.Profile) error {
	if err := CheckSocialTypes(r.db, &profile); err != nil {
		return err
	}

	profile.GamesListed = 0

	res := r.db.Create(&profile)
	if res.Error != nil {
		return utilErrs.FromGORM(res, "failed to create profile")
	}

	return nil
}

func (r *gameListRepository) SaveProfile(profile entity.Profile) error {
	if err := CheckSocialTypes(r.db, &profile); err != nil {
		return err
	}

	res := r.db.Save(&profile)
	if res.Error != nil {
		return utilErrs.FromGORM(res, "failed to save profile")
	}

	return nil
}

func (r *gameListRepository) GetAllProfiles() ([]entity.ProfileInfo, error) {
	var profiles []entity.ProfileInfo
	res := r.db.Model(&entity.Profile{}).Find(&profiles)
	if res.Error != nil {
		return nil, utilErrs.FromGORM(res, "failed to get profiles")
	}

	return profiles, nil
}

func (r *gameListRepository) GetProfile(login entity.ProfileCreds) (*entity.Profile, error) {
	var profile entity.Profile
	res := r.db.First(&profile, login)
	if res.Error != nil {
		return nil, utilErrs.FromGORM(res, "failed to get profile")
	}

	return &profile, nil
}

func (r *gameListRepository) SaveRefreshToken(nickname string, tokenString string) error {
	userID, err := r.findUserIDByNickname(nickname)
	if err != nil {
		return err
	}

	refreshToken := entity.RefreshToken{
		ProfileID: userID,
		Token:     tokenString,
	}

	res := r.db.Create(&refreshToken)
	if res.Error != nil {
		return utilErrs.FromGORM(res, "unable to save the refresh token")
	}

	return nil
}

func (r *gameListRepository) FindRefreshToken(nickname string, tokenString string) error {
	var result entity.RefreshToken
	res := r.db.Table("refresh_token, profile").Select("refresh_token.token").Where(
		"refresh_token.token = ?", tokenString).Where(
		"refresh_token.profile_id = profile.id").Where(
		"profile.nickname = ?", nickname).Where(
		"refresh_token.deleted_at IS NULL").
		Limit(1).
		Scan(&result)

	if res.Error != nil || res.RowsAffected == 0 {
		return utilErrs.FromGORM(res, "failed to find refresh token")
	}

	return nil
}

func (r *gameListRepository) DeleteRefreshToken(tokenString string) error {
	res := r.db.Where("token = ?", tokenString).Delete(&entity.RefreshToken{})
	if res.Error != nil {
		return utilErrs.FromGORM(res, "failed to delete refresh token")
	}

	return nil
}

func (r *gameListRepository) DeleteAllUserRefreshTokens(nickname string) error {
	userID, err := r.findUserIDByNickname(nickname)
	if err != nil {
		return err
	}

	res := r.db.Where("profile_id = ?", userID).Delete(&entity.RefreshToken{})
	if res.Error != nil {
		return utilErrs.FromGORM(res, "failed to delete refresh token")
	}

	return nil
}

func (r *gameListRepository) SaveSocialType(socialType entity.SocialType) error {
	res := r.db.Save(&socialType)
	if res.Error != nil {
		return utilErrs.FromGORM(res, "failed to save social type")
	}

	return nil
}

func (r *gameListRepository) GetAllSocialTypes() ([]entity.SocialType, error) {
	var socialTypes []entity.SocialType

	res := r.db.Find(&socialTypes)
	if res.Error != nil {
		return nil, utilErrs.FromGORM(res, "failed to find social types")
	}

	return socialTypes, nil
}

func (r *gameListRepository) findUserIDByNickname(nickname string) (uint64, error) {
	var userID uint64
	res := r.db.Table("profile").Select("id").Take(&userID, map[string]string{"nickname": nickname})
	if res.Error != nil {
		return 0, utilErrs.FromGORM(res, fmt.Sprintf("failed to find user with nickname \"%s\"", nickname))
	}

	return userID, nil
}
