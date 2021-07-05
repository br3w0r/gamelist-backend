package service

import (
	"context"
	"io"
	"log"
	"time"

	"bitbucket.org/br3w0r/gamelist-backend/entity"
	pb "bitbucket.org/br3w0r/gamelist-backend/proto"
	"bitbucket.org/br3w0r/gamelist-backend/repository"
	"google.golang.org/grpc"
)

type GameListService interface {
	SaveGame(game entity.GameProperties) error
	GetAllGames() []entity.GameProperties

	SaveGenre(genre entity.Genre) error
	GetAllGenres() []entity.Genre

	SavePlatform(platform entity.Platform) error
	GetAllPlatforms() []entity.Platform

	CreateProfile(profile entity.ProfileInfo) error
	SaveProfile(profile entity.ProfileInfo) error
	GetAllProfiles() []entity.ProfileInfo

	SaveSocialType(socialType entity.SocialType) error
	GetAllSocialTypes() []entity.SocialType

	// gRPC
	ScrapeGames()
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

func (s *gameListService) CreateProfile(profile entity.ProfileInfo) error {
	return s.repo.CreateProfile(profile)
}

func (s *gameListService) SaveProfile(profile entity.ProfileInfo) error {
	return s.repo.SaveProfile(profile)
}

func (s *gameListService) GetAllProfiles() []entity.ProfileInfo {
	return s.repo.GetAllProfiles()
}

func (s *gameListService) SaveSocialType(socialType entity.SocialType) error {
	return s.repo.SaveSocialType(socialType)
}

func (s *gameListService) GetAllSocialTypes() []entity.SocialType {
	return s.repo.GetAllSocialTypes()
}

// gRPC
func (s *gameListService) ScrapeGames() {
	// Connection
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}

	conn, err := grpc.Dial("localhost:8888", opts...)
	if err != nil {
		log.Printf("failed to dial: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewGameScrapeClient(conn)

	// Getting games from scraper
	stream, err := client.ScrapeGames(context.Background(), &pb.Empty{})
	if err != nil {
		log.Printf("<ScrapeGames>: failed to set up stream: %v", err)
		return
	}

	counter := 1
	t := time.Now()
	for {
		log.Printf("Adding game: %d", counter)
		game, err := stream.Recv()

		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("rpc error: %v", err)
			return
		}

		s.repo.SaveGame(game.ConvertToEntity())
		counter++
	}
	log.Printf("Time elapsed: %v", time.Since(t))
}
