package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"said-and-dot-backend/internal/common/config"
	"said-and-dot-backend/internal/common/logger"
	"said-and-dot-backend/internal/database"
	"said-and-dot-backend/internal/server"
)

var (
	dotEnvConfigFilePath string
	preforkStatus        bool
)

func init() {
	flag.StringVar(&dotEnvConfigFilePath, "config", "./configs/.env", "The application configurations")
	flag.BoolVar(&preforkStatus, "prefork", false, "Run the app in Prefork mode [multiple Go processes]")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := config.Load(dotEnvConfigFilePath); err != nil {
		log.Fatal(err)
	}

	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.GetString("DB_USER", "postgres"),
		config.GetString("DB_PASSWORD", "postgres"),
		config.GetString("DB_HOST", "127.0.0.1"),
		config.GetInt("DB_PORT", 5432),
		config.GetString("DB_DATABASE", "snd-data"),
	)

	println(config.GetString("DB_HOST", "127.0.0.1"))

	println(dbUrl)

	dbInstance, err := database.New(ctx, dbUrl)
	if err != nil {
		log.Println(dbUrl)
		log.Fatal(err)
	}

	serverInstance := server.New(
		dbInstance,
		logger.NewLogger(config.GetBool("DEBUG", false)),
		&server.Config{
			AppName: config.GetString("APP_NAME", "monoBackend"),
			Host:    config.GetString("APP_HOST", "127.0.0.1"),
			Port:    config.GetString("APP_PORT", "5000"),
		},
		fiber.Config{
			Prefork: preforkStatus,
		},
	)

	serverInstance.Listen()
}
