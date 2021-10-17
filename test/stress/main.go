package test

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/br3w0r/gamelist-backend/entity"
	"github.com/br3w0r/gamelist-backend/repository"
)

func parseEntries(op []string) (uint64, error) {
	if len(op) < 2 {
		return 0, fmt.Errorf("wrong option: %s. Amount of entries to get must be specified", op[0])
	}
	n, err := strconv.ParseUint(op[1], 0, 0)
	if err != nil {
		return 0, fmt.Errorf("wrong option: %s. Amount of entries to get must be uint", op[0])
	}
	return n, nil
}

func getGame(repo repository.GamelistRepository, n uint64) {
	var i uint64
	for i = 1; i <= n; i++ {
		repo.GetGameDetails("test", i)
	}
}

func userCreation(repo repository.GamelistRepository, n uint64) {
	var i uint64
	for i = 1; i <= n; i++ {
		repo.CreateProfile(entity.Profile{
			ProfileInfo: entity.ProfileInfo{
				Nickname: fmt.Sprintf("test%d", i),
			},
			Email:    fmt.Sprintf("test%d@mail.com", i),
			Password: "123",
		})
	}
}

func findUser(repo repository.GamelistRepository, n uint64) {
	var i uint64
	for i = 1; i <= n; i++ {
		repo.GetProfile(entity.ProfileCreds{Nickname: fmt.Sprintf("test%d", i)})
	}
}

func findGames(repo repository.GamelistRepository) {
	games, _ := repo.GetAllGames()

	for _, game := range games {
		name := game.Name[:len(game.Name)/2]
		repo.SearchGames(name)
	}
}

func RunStress(repo repository.GamelistRepository, options []string) {
	repo.CreateProfile(entity.Profile{
		ProfileInfo: entity.ProfileInfo{
			Nickname: "test",
		},
		Email:    "test@mail.com",
		Password: "123",
	})

	for _, i := range options {
		log.Printf(">>> Test: %s", i)
		t := time.Now()
		op := strings.Split(i, "=")
		switch op[0] {
		case "get_user_games", "get_all_games": // These tests will be implemented soon
			log.Printf("%s unimplemented", op[0])
		case "user_creation":
			n, err := parseEntries(op)
			if err != nil {
				log.Println(err.Error())
			} else {
				userCreation(repo, n)
			}
		case "find_user":
			n, err := parseEntries(op)
			if err != nil {
				log.Println(err.Error())
			} else {
				findUser(repo, n)
			}
		case "get_game":
			n, err := parseEntries(op)
			if err != nil {
				log.Println(err.Error())
			} else {
				getGame(repo, n)
			}
		case "find_games":
			findGames(repo)
		default:
			log.Printf("%s unimplemented", op[0])
		}
		log.Printf("<<< Test ended. Time elapsed: %v", time.Since(t))
	}
}
