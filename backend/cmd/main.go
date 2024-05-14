package main

import (
	"flag"
	"os"

	"github.com/ShareCampus/RecRem/backend/pkg/database"
	"github.com/ShareCampus/RecRem/backend/pkg/recrem"
	"github.com/ShareCampus/RecRem/backend/pkg/utils/logger"
	"github.com/joho/godotenv"
)

const HTTP_LISTEN_PORT_STR = "HTTP_LISTEN_PORT"

func main() {
	var envFilePath string
	flag.StringVar(&envFilePath, "env", "", "Path to .env file")
	flag.Parse()
	if envFilePath != "" {
		logger.Info("load env")
		err := godotenv.Load(envFilePath)
		if err != nil {
			logger.Info("Failed to load .env file, using system ones.")
		} else {
			logger.Infof("Loaded .env file from %s", envFilePath)
		}
	}

	listenPort := os.Getenv(HTTP_LISTEN_PORT_STR)
	db := database.SetupDB(nil)
	server := recrem.NewServer(
		recrem.ListenPort(listenPort),
		recrem.UseDB(db),
	)
	panic(server)
}
