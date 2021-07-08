package repository

import (
	"fmt"
	"os"

	"bitbucket.org/br3w0r/gamelist-backend/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GamelistRepository interface {
	SaveGame(game entity.GameProperties) error
	GetAllGames() []entity.GameProperties

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

func NewGamelistRepository(dbName string, forceMigrate bool) GamelistRepository {
	var db *gorm.DB
	_, err := os.Stat(dbName)

	if os.IsNotExist(err) || forceMigrate {
		db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
		if err != nil {
			panic("Failed to connect database.")
		}
		db.AutoMigrate(&entity.GameProperties{}, &entity.Genre{},
			&entity.Platform{}, &entity.Profile{}, &entity.RefreshToken{}, &entity.Social{},
			&entity.SocialType{}, &entity.ProfileGames{}, &entity.ListType{})
	} else {
		db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
		if err != nil {
			panic("Failed to connect database.")
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
	err := r.db.Table("refresh_tokens, profiles").Select("refresh_tokens.token").Where(
		"refresh_tokens.token = ?", tokenString).Where(
		"refresh_tokens.profile_id = profiles.id").Where(
		"profiles.nickname = ?", nickname).Where(
		"refresh_tokens.deleted_at IS NULL").Scan(&result).Error

	if err != nil {
		return err
	}

	fmt.Println(len(result), result)

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
	err := r.db.Table("profiles").Select("id").Take(&userID, map[string]string{"nickname": nickname}).Error
	if err != nil {
		return 0, fmt.Errorf("unable to find user with given nickname: %v", err)
	}
	return userID, nil
}
