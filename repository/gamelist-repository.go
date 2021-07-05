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
	GetProfile(nickname string) (*entity.Profile, error)

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
			&entity.Platform{}, &entity.Profile{}, &entity.Social{},
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
	for i := range profile.Socials {
		err := r.db.First(&entity.SocialType{}, profile.Socials[i].TypeID).Error
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("couldn't find social of type %d", profile.Socials[i].TypeID)
		}
		if err != nil {
			return err
		}
	}

	profile.GamesListed = 0

	return r.db.Create(&profile).Error
}

func (r *gameListRepository) SaveProfile(profile entity.Profile) error {
	return r.db.Save(&profile).Error
}

func (r *gameListRepository) GetAllProfiles() []entity.ProfileInfo {
	var profiles []entity.ProfileInfo
	r.db.Model(&entity.Profile{}).Find(&profiles)
	return profiles
}

func (r *gameListRepository) GetProfile(nickname string) (*entity.Profile, error) {
	var profile entity.Profile
	err := r.db.First(&profile, map[string]string{"nickname": nickname}).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *gameListRepository) SaveSocialType(socialType entity.SocialType) error {
	return r.db.Save(&socialType).Error
}

func (r *gameListRepository) GetAllSocialTypes() []entity.SocialType {
	var socialTypes []entity.SocialType
	r.db.Find(&socialTypes)
	return socialTypes
}
