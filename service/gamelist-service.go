package service

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/br3w0r/gamelist-backend/entity"
	pb "github.com/br3w0r/gamelist-backend/proto"
	"github.com/br3w0r/gamelist-backend/repository"
	utilErrs "github.com/br3w0r/gamelist-backend/util/errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type GameListService interface {
	SaveGame(game entity.GameProperties) error
	GetAllGames() ([]entity.GameProperties, error)
	GetAllGamesTyped(nickname string, last uint64, batchSize int) ([]entity.TypedGameListProperties, error)
	GetUserGameList(nickname string) ([]entity.TypedGameListProperties, error)
	SearchGames(name string) ([]entity.GameSearchResult, error)
	GetGameDetails(nickname string, gameId uint64) (*entity.GameDetailsResponse, error)

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
	CheckLogin(login entity.LoginProfile) (*entity.Profile, error)

	SaveSocialType(socialType entity.SocialType) error
	GetAllSocialTypes() ([]entity.SocialType, error)

	// gRPC
	ScrapeGames()
}

type gameListService struct {
	repo               repository.GamelistRepository
	scraperGRPCAddress string
}

func NewGameListService(repo repository.GamelistRepository, scraperGRPCAddress string) GameListService {
	return &gameListService{repo, scraperGRPCAddress}
}

func (s *gameListService) SaveGame(game entity.GameProperties) error {
	return s.repo.SaveGame(game)
}

func (s *gameListService) GetAllGames() ([]entity.GameProperties, error) {
	return s.repo.GetAllGames()
}

func (s *gameListService) GetAllGamesTyped(nickname string, last uint64, batchSize int) ([]entity.TypedGameListProperties, error) {
	return s.repo.GetAllGamesTyped(nickname, last, batchSize)
}

func (s *gameListService) GetUserGameList(nickname string) ([]entity.TypedGameListProperties, error) {
	return s.repo.GetUserGameList(nickname)
}

func (s *gameListService) SearchGames(name string) ([]entity.GameSearchResult, error) {
	return s.repo.SearchGames(name)
}

func (s *gameListService) GetGameDetails(nickname string, gameId uint64) (*entity.GameDetailsResponse, error) {
	return s.repo.GetGameDetails(nickname, gameId)
}

func (s *gameListService) CreateListType(listType entity.ListType) error {
	return s.repo.CreateListType(listType)
}

func (s *gameListService) GetAllListTypes() ([]entity.ListType, error) {
	return s.repo.GetAllListTypes()
}

func (s *gameListService) ListGame(nickname string, gameId uint64, listType uint64) error {
	return s.repo.ListGame(nickname, gameId, listType)
}

func (s *gameListService) SaveGenre(genre entity.Genre) error {
	return s.repo.SaveGenre(genre)
}

func (s *gameListService) GetAllGenres() ([]entity.Genre, error) {
	return s.repo.GetAllGenres()
}

func (s *gameListService) SavePlatform(platform entity.Platform) error {
	return s.repo.SavePlatform(platform)
}

func (s *gameListService) GetAllPlatforms() ([]entity.Platform, error) {
	return s.repo.GetAllPlatforms()
}

func (s *gameListService) CreateProfile(profile entity.Profile) error {
	// Encrypting password
	hash, err := bcrypt.GenerateFromPassword([]byte(profile.Password), 10)
	if err != nil {
		return err
	}

	profile.Password = string(hash)

	return s.repo.CreateProfile(profile)
}

func (s *gameListService) SaveProfile(profile entity.Profile) error {
	// Generate hash for new password if it was changed
	if len(profile.Password) > 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(profile.Password), 10)
		if err != nil {
			return err
		}

		profile.Password = string(hash)
	}

	return s.repo.SaveProfile(profile)
}

func (s *gameListService) GetAllProfiles() ([]entity.ProfileInfo, error) {
	return s.repo.GetAllProfiles()
}

func (s *gameListService) CheckLogin(login entity.LoginProfile) (*entity.Profile, error) {
	profile, err := s.repo.GetProfile(entity.ProfileCreds{
		Nickname: login.Nickname,
		Email:    login.Email,
	})
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(profile.Password), []byte(login.Password))
	if err != nil {
		return nil, utilErrs.New(utilErrs.Unauthorized, err, "incorrect password")
	}

	return profile, nil
}

func (s *gameListService) SaveSocialType(socialType entity.SocialType) error {
	return s.repo.SaveSocialType(socialType)
}

func (s *gameListService) GetAllSocialTypes() ([]entity.SocialType, error) {
	return s.repo.GetAllSocialTypes()
}

// gRPC
//nolint Will be remade
func (s *gameListService) ScrapeGames() {
	// Connection
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}

	conn, err := grpc.Dial(s.scraperGRPCAddress+":8888", opts...)
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
