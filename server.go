package main

import (
	"strings"

	"github.com/br3w0r/gamelist-backend/helpers"
	"github.com/br3w0r/gamelist-backend/server"
)

var (
	PRODUCTION_MODE      string = helpers.GetEnvOrDefault("PRODUCTION_MODE", "0")
	SERVE_STATIC         string = helpers.GetEnvOrDefault("SERVE_STATIC", "1")
	FORCE_MIGRATE        string = helpers.GetEnvOrDefault("FORCE_MIGRATE", "0")
	FORCE_SCRAPE         string = helpers.GetEnvOrDefault("FORCE_SCRAPE", "0")
	STATIC_DIR           string = helpers.GetEnvOrDefault("STATIC_FOLDER", "../gamelist-frontend/gamelist/dist")
	DATABASE_DIR         string = helpers.GetEnvOrDefault("DATABASE_DIR", ".")
	SCRAPER_GRPC_ADDRESS string = helpers.GetEnvOrDefault("SCRAPER_GRPC_ADDRESS", "localhost")
	STRESS_TEST          string = helpers.GetEnvOrDefault("STRESS_TEST", "0")
	STRESS_TEST_OPTIONS  string = helpers.GetEnvOrDefault("STRESS_TEST_OPTIONS", "user_creation,get_game=75,get_all_games,get_user_games")
)

func main() {
	options := server.ServerOptions{
		Production:         PRODUCTION_MODE == "1",
		ServeStatic:        SERVE_STATIC == "1",
		ForceMigrate:       FORCE_MIGRATE == "1",
		ForceScrape:        FORCE_SCRAPE == "1",
		StaticDir:          STATIC_DIR,
		DatabaseDir:        DATABASE_DIR,
		ScraperGRPCAddress: SCRAPER_GRPC_ADDRESS,
		StressTest:         STRESS_TEST == "1",
		StressTestOptions:  strings.Split(STRESS_TEST_OPTIONS, ","),
	}

	server := server.NewServer(options)

	server.Run(":8080")
}
