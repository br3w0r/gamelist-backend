package service

import (
	"bitbucket.org/br3w0r/gamelist-backend/entity"
	"bitbucket.org/br3w0r/gamelist-backend/repository"
)

type GameListService interface {
	SaveGame(game entity.GameProperties) error
	GetAllGames() []entity.GameProperties

	SaveGenre(genre entity.Genre) error
	GetAllGenres() []entity.Genre

	SavePlatform(platform entity.Platform) error
	GetAllPlatforms() []entity.Platform
}

type gameListService struct {
	repo repository.GamelistRepository
}

func NewGameListService(repo repository.GamelistRepository) GameListService {
	return &gameListService{repo}
}

func (s *gameListService) SaveGame(game entity.GameProperties) error {
	return s.repo.SaveGame(game)
}

func (s *gameListService) GetAllGames() []entity.GameProperties {
	return s.repo.GetAllGames()
}

func (s *gameListService) SaveGenre(genre entity.Genre) error {
	return s.repo.SaveGenre(genre)
}

func (s *gameListService) GetAllGenres() []entity.Genre {
	return s.repo.GetAllGenres()
}

func (s *gameListService) SavePlatform(platform entity.Platform) error {
	return s.repo.SavePlatform(platform)
}

func (s *gameListService) GetAllPlatforms() []entity.Platform {
	return s.repo.GetAllPlatforms()
}
