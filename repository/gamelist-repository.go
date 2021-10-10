package repository

import (
	"fmt"
	"log"

	"github.com/br3w0r/gamelist-backend/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GamelistRepository interface {
	SaveGame(game entity.GameProperties) error
	GetAllGames() []entity.GameProperties
	GetAllGamesTyped(nickname string, last uint64, batchSize int) []entity.TypedGameListProperties
	GetUserGameList(nickname string) []entity.TypedGameListProperties
	SearchGames(name string) []entity.GameSearchResult
	GetGameDetails(nickname string, id uint64) (*entity.GameDetailsResponse, error)

	CreateListType(listType entity.ListType) error
	GetAllListTypes() []entity.ListType
	ListGame(nickname string, gameId uint64, listType uint64) error

	SaveGenre(genre entity.Genre) error
	GetAllGenres() []entity.Genre

	SavePlatform(platform entity.Platform) error
	GetAllPlatforms() []entity.Platform

	CreateProfile(profile entity.Profile) error
	SaveProfile(profile entity.Profile) error
	GetAllProfiles() []entity.ProfileInfo
	GetProfile(login entity.ProfileCreds) (*entity.Profile, error)

	SaveRefreshToken(nickname string, tokenString string) error
	FindRefreshToken(nickname string, tokenString string) error
	DeleteRefreshToken(tokenString string) error
	DeleteAllUserRefreshTokens(nickname string) error

	SaveSocialType(socialType entity.SocialType) error
	GetAllSocialTypes() []entity.SocialType
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

func NewGamelistRepository(dbName string, forceMigrate bool, dialector gorm.Dialector) GamelistRepository {
	var db *gorm.DB
	var err error

	if forceMigrate {
		db, err = gorm.Open(dialector, &gorm.Config{})
		if err != nil {
			panic(ErrDbConnection)
		}
		db.AutoMigrate(&entity.GameProperties{}, &entity.Genre{},
			&entity.Platform{}, &entity.Profile{}, &entity.RefreshToken{}, &entity.ProfileGame{},
			&entity.Social{}, &entity.SocialType{}, &entity.ListType{})

		err := db.Model(&entity.ListType{}).First(nil).Error
		if err == gorm.ErrRecordNotFound {
			log.Println("Creating default list types...")
			listTypes := []entity.ListType{
				{Name: "Played"},
				{Name: "Playing"},
				{Name: "Want to play"},
			}

			db.Create(&listTypes)
		} else if err != nil {
			panic("Failed to get first list type")
		}
	} else {
		db, err = gorm.Open(dialector, &gorm.Config{})
		if err != nil {
			panic(ErrDbConnection)
		}
	}

	return &gameListRepository{
		db: db,
	}
}

func (r *gameListRepository) SaveGame(game entity.GameProperties) error {
	for i := range game.Platforms {
		err := r.db.First(&game.Platforms[i], game.Platforms[i]).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
	}
	for i := range game.Genres {
		err := r.db.First(&game.Genres[i], game.Genres[i]).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
	}
	return r.db.Save(&game).Error
}

func (r *gameListRepository) GetAllGames() []entity.GameProperties {
	var games []entity.GameProperties
	r.db.Preload(clause.Associations).Find(&games)
	return games
}

func (r *gameListRepository) GetAllGamesTyped(nickname string, last uint64, batchSize int) []entity.TypedGameListProperties {
	if batchSize > GAMES_BATCH_SIZE_LIMIT {
		batchSize = GAMES_BATCH_SIZE_LIMIT
	}

	userId, err := r.findUserIDByNickname(nickname)
	if err != nil {
		return nil
	}
	var games []entity.TypedGameListProperties
	r.db.Table("game_properties").
		Joins("left join profile_game on game_properties.id = profile_game.game_id and profile_game.profile_id = ?", userId).
		Where("game_properties.id > ?", last).
		Limit(batchSize).
		Scan(&games)

	return games
}

func (r *gameListRepository) GetUserGameList(nickname string) []entity.TypedGameListProperties {
	var games []entity.TypedGameListProperties
	r.db.Table("game_properties").Select(
		"game_properties.id, game_properties.name, game_properties.image_url, game_properties.year_released, profile_game.list_type_id",
	).Joins(
		"join profile_game on game_properties.id = profile_game.game_id and profile_game.list_type_id != 0",
	).Joins(
		"join profile on profile_game.profile_id = profile.id and profile.nickname = ?",
		nickname,
	).Scan(&games)
	return games
}

func (r *gameListRepository) SearchGames(name string) []entity.GameSearchResult {
	var games []entity.GameSearchResult
	if len(name) > 1 {
		r.db.Table("game_properties").Where("name LIKE ?", name+"%").Limit(10).Find(&games)
	}
	return games
}

func (r *gameListRepository) GetGameDetails(nickname string, gameId uint64) (*entity.GameDetailsResponse, error) {
	userId, err := r.findUserIDByNickname(nickname)
	if err != nil {
		return nil, err
	}
	var gameDetails entity.GameDetailsResponse
	err = r.db.Table("game_properties").
		Joins("left join profile_game on game_properties.id = profile_game.game_id and profile_game.profile_id = ?", userId).
		Where("game_properties.id = ?", gameId).
		Scan(&(gameDetails.Game)).Error

	if gameDetails.Game.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	err = r.db.Table("platform").Select("platform.name").
		Joins("inner join game_platforms on game_platforms.game_properties_id = ? and game_platforms.platform_id = platform.id", gameId).
		Scan(&(gameDetails.Platforms)).Error

	if err != nil {
		return nil, err
	}

	err = r.db.Table("genre").Select("genres.name").
		Joins("inner join game_genres on game_genres.game_properties_id = ? and game_genres.genre_id = genre.id", gameId).
		Scan(&(gameDetails.Genres)).Error

	if err != nil {
		return nil, err
	}

	return &gameDetails, nil
}

func (r *gameListRepository) CreateListType(listType entity.ListType) error {
	return r.db.Create(&listType).Error
}

func (r *gameListRepository) GetAllListTypes() []entity.ListType {
	var types []entity.ListType
	r.db.Find(&types)
	return types
}

func (r *gameListRepository) ListGame(nickname string, gameId uint64, listType uint64) error {
	userId, err := r.findUserIDByNickname(nickname)
	if err != nil {
		return err
	}

	err = r.db.First(&entity.GameProperties{}, gameId).Error
	if err == gorm.ErrRecordNotFound {
		return fmt.Errorf("couldn't find game with id: %d", gameId)
	}
	if err != nil {
		return err
	}

	if listType != 0 {
		err = r.db.First(&entity.ListType{}, listType).Error
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("couldn't find list type with id: %d", listType)
		}
		if err != nil {
			return err
		}
	}

	listGame := entity.ProfileGame{
		ProfileID:  userId,
		GameID:     gameId,
		ListTypeID: listType,
	}

	return r.db.Save(&listGame).Error
}

func (r *gameListRepository) SaveGenre(genre entity.Genre) error {
	return r.db.Save(&genre).Error
}

func (r *gameListRepository) GetAllGenres() []entity.Genre {
	var genres []entity.Genre
	r.db.Find(&genres)
	return genres
}

func (r *gameListRepository) SavePlatform(platform entity.Platform) error {
	return r.db.Save(&platform).Error
}

func (r *gameListRepository) GetAllPlatforms() []entity.Platform {
	var platforms []entity.Platform
	r.db.Find(&platforms)
	return platforms
}

func (r *gameListRepository) CreateProfile(profile entity.Profile) error {
	if err := CheckSocialTypes(r.db, &profile); err != nil {
		return err
	}

	profile.GamesListed = 0

	return r.db.Create(&profile).Error
}

func (r *gameListRepository) SaveProfile(profile entity.Profile) error {
	if err := CheckSocialTypes(r.db, &profile); err != nil {
		return err
	}

	return r.db.Save(&profile).Error
}

func (r *gameListRepository) GetAllProfiles() []entity.ProfileInfo {
	var profiles []entity.ProfileInfo
	r.db.Model(&entity.Profile{}).Find(&profiles)
	return profiles
}

func (r *gameListRepository) GetProfile(login entity.ProfileCreds) (*entity.Profile, error) {
	var profile entity.Profile
	err := r.db.First(&profile, login).Error
	if err != nil {
		return nil, err
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

	err = r.db.Create(&refreshToken).Error
	if err != nil {
		return fmt.Errorf("unable to save the refresh token: %v", err)
	}
	return nil
}

func (r *gameListRepository) FindRefreshToken(nickname string, tokenString string) error {
	var result []entity.RefreshToken
	err := r.db.Table("refresh_token, profile").Select("refresh_token.token").Where(
		"refresh_token.token = ?", tokenString).Where(
		"refresh_token.profile_id = profile.id").Where(
		"profile.nickname = ?", nickname).Where(
		"refresh_token.deleted_at IS NULL").Scan(&result).Error

	if err != nil {
		return err
	}

	if len(result) == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *gameListRepository) DeleteRefreshToken(tokenString string) error {
	return r.db.Where("token = ?", tokenString).Delete(&entity.RefreshToken{}).Error
}

func (r *gameListRepository) DeleteAllUserRefreshTokens(nickname string) error {
	userID, err := r.findUserIDByNickname(nickname)
	if err != nil {
		return err
	}

	return r.db.Where("profile_id = ?", userID).Delete(&entity.RefreshToken{}).Error
}

func (r *gameListRepository) SaveSocialType(socialType entity.SocialType) error {
	return r.db.Save(&socialType).Error
}

func (r *gameListRepository) GetAllSocialTypes() []entity.SocialType {
	var socialTypes []entity.SocialType
	r.db.Find(&socialTypes)
	return socialTypes
}

func (r *gameListRepository) findUserIDByNickname(nickname string) (uint64, error) {
	var userID uint64
	err := r.db.Table("profile").Select("id").Take(&userID, map[string]string{"nickname": nickname}).Error
	if err != nil {
		return 0, fmt.Errorf("unable to find user with given nickname: %v", err)
	}
	return userID, nil
}
