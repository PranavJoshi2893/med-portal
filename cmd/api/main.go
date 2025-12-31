package main

import (
	"log"

	"github.com/PranavJoshi2893/boilerplate/internal/config"
	"github.com/PranavJoshi2893/boilerplate/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	routes := server.Routes()

	srv := server.NewServer(cfg, routes)

	log.Println("server is running on port", cfg.ServerPort)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Shutdown")

}
