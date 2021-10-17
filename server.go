package main

import (
	"log"
	"strings"

	"github.com/br3w0r/gamelist-backend/helpers"
	"github.com/br3w0r/gamelist-backend/repository"
	"github.com/br3w0r/gamelist-backend/server"
)

var (
	PORT                 string = helpers.GetEnvOrDefault("PORT", "8080")
	PRODUCTION_MODE      string = helpers.GetEnvOrDefault("PRODUCTION_MODE", "0")
	SERVE_STATIC         string = helpers.GetEnvOrDefault("SERVE_STATIC", "1")
	FORCE_SCRAPE         string = helpers.GetEnvOrDefault("FORCE_SCRAPE", "0")
	STATIC_DIR           string = helpers.GetEnvOrDefault("STATIC_FOLDER", "../gamelist-frontend/gamelist/dist")
	DATABASE_DIST        string = helpers.GetEnvOrDefault("DATABASE_DIST", "./gamelist.db")
	SCRAPER_GRPC_ADDRESS string = helpers.GetEnvOrDefault("SCRAPER_GRPC_ADDRESS", "localhost")
	STRESS_TEST          string = helpers.GetEnvOrDefault("STRESS_TEST", "0")
	STRESS_TEST_OPTIONS  string = helpers.GetEnvOrDefault("STRESS_TEST_OPTIONS", "user_creation,get_game=75,get_all_games,get_user_games")
	DB_HOST              string = helpers.GetEnvOrDefault("DB_HOST", "localhost")
	DB_PORT              string = helpers.GetEnvOrDefault("DB_PORT", "5432")
	DB_USER              string = helpers.GetEnvOrDefault("DB_USER", "postgres")
	DB_NAME              string = helpers.GetEnvOrDefault("DB_NAME", "gamelist")
	DB_PASSWORD          string = helpers.GetEnvOrDefault("DB_PASSWORD", "pgpass")
	DB_SSL               string = helpers.GetEnvOrDefault("DB_SSL", "0")
	DB_TIMEZONE          string = helpers.GetEnvOrDefault("DB_TIMEZONE", "")
)

func main() {
	var scraperAsync bool
	if STRESS_TEST == "1" {
		scraperAsync = true
	} else {
		scraperAsync = false
	}
	options := server.ServerOptions{
		Production:         PRODUCTION_MODE == "1",
		ServeStatic:        SERVE_STATIC == "1",
		ForceScrape:        FORCE_SCRAPE == "1",
		StaticDir:          STATIC_DIR,
		DatabaseDist:       DATABASE_DIST,
		ScraperGRPCAddress: SCRAPER_GRPC_ADDRESS,
		ScraperAsync:       scraperAsync,
		StressTest:         STRESS_TEST == "1",
		StressTestOptions:  strings.Split(STRESS_TEST_OPTIONS, ","),
		SilentMode:         false,
		DBConfig: &repository.DBConfig{
			Host:     DB_HOST,
			Port:     DB_PORT,
			User:     DB_USER,
			DBName:   DB_NAME,
			Password: DB_PASSWORD,
			SSL:      DB_SSL == "1",
			TimeZone: DB_TIMEZONE,
		},
	}

	server := server.NewServer(options)

	if !options.StressTest {
		err := server.Run(":" + PORT)
		if err != nil {
			log.Fatalf("failed to start server: %s", err)
		}
	}
}
