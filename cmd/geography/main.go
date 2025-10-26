package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/NarzhanProduction/Geography/cmd/config"
	api_http "github.com/NarzhanProduction/Geography/internal/api/http"
	"github.com/NarzhanProduction/Geography/internal/api/http/handlers"
	"github.com/NarzhanProduction/Geography/internal/database"
	"github.com/NarzhanProduction/Geography/internal/database/repository"
	"github.com/NarzhanProduction/Geography/internal/pkg/logger"
	"github.com/NarzhanProduction/Geography/internal/service"
)

func main() {
	mainLogger := logger.New("geography")

	if err := run(mainLogger); err != nil {
		log.Fatalf("error happened: %v", err)
	}
}

func run(mainLogger logger.Logger) error {
	ctx := context.Background()

	cfg := config.InitConfig()

	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)

	repo, err := database.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("error starting database: %v", err)
	}
	userRepo := repository.NewUserRepository(repo.DB)
	jwtTtl, _ := time.ParseDuration(cfg.JWTttl)
	usrSrv := service.NewUserService(userRepo, cfg.JWTKey, jwtTtl)
	gptSrv := service.NewChatbotService(cfg.OpenAPIKey, mainLogger)
	srv := service.NewService(usrSrv, gptSrv)
	handlers := handlers.NewHandler(srv)

	router := api_http.InitRouter()
	api_http.RegisterHandlers(router, handlers)
	router.Start(fmt.Sprintf("%s:%s", cfg.Address, cfg.Port))
	return nil
}
