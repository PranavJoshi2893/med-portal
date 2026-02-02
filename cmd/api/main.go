package main

import (
	"log"

	"github.com/PranavJoshi2893/med-portal/internal/config"
	"github.com/PranavJoshi2893/med-portal/internal/database"
	"github.com/PranavJoshi2893/med-portal/internal/handler"
	"github.com/PranavJoshi2893/med-portal/internal/repository"
	"github.com/PranavJoshi2893/med-portal/internal/server"
	"github.com/PranavJoshi2893/med-portal/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, cfg.Pepper, cfg.AccessTokenKey, cfg.RefreshTokenKey)
	authHandler := handler.NewAuthHandler(authService)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	routes := server.Routes(authHandler, userHandler, cfg)

	srv := server.NewServer(cfg, db, routes)

	log.Println("server is running on port", cfg.ServerPort)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Shutdown")

}
