package main

import (
	"log"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/config"
	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	log.Println("Server is listening on ", config.Params.ServerAddress)
	server.Configure()
	if err := server.Run(); err != nil {
		log.Fatal("Error starting server")
	}
}
