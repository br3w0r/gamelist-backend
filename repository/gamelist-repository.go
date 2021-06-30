package repository

import (
	"encoding/json"
	"errors"
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
			&entity.Platform{}, &entity.ProfileInfo{}, &entity.Social{},
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
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				out, _ := json.Marshal(&game.Platforms[i])
				return errors.New("Can't find platform: " + string(out))
			}
			return err
		}
	}
	for i := range game.Genres {
		err := r.db.First(&game.Genres[i], game.Genres[i]).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				out, _ := json.Marshal(&game.Genres[i])
				return errors.New("Can't find genre: " + string(out))
			}
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
