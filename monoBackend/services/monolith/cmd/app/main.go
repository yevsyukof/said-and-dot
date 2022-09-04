package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"log"
	"said-and-dot-backend/pkg/config"
	"said-and-dot-backend/pkg/logger"
	"said-and-dot-backend/pkg/store/postgres"
	"said-and-dot-backend/services/monolith/internal/server"
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
		log.Println("-- dotEnvFile not found --")
		log.Fatal(err)
	}

	store, err := postgres.NewStore(ctx, &postgres.Config{
		Host:     config.GetString("DB_HOST", "127.0.0.1"),
		Port:     config.GetUint16("DB_PORT", 5432),
		Username: config.GetString("DB_USER", "postgres"),
		Password: config.GetString("DB_PASSWORD", "postgres"),
		DBName:   config.GetString("DB_DATABASE", "snd-data"),
	})
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	serv := server.New(
		store,
		logger.NewLogger(config.GetBool("DEBUG", false)),
		&server.Config{
			AppName: config.GetString("APP_NAME", "monoBackend"),
			Host:    config.GetString("APP_HOST", ""),
			Port:    config.GetString("APP_PORT", "5000"),
		},
		fiber.Config{
			Prefork: preforkStatus,
		},
	)

	serv.Listen()
}
