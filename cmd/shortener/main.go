package main

import (
	"log"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/server"
)

func main() {
	log.Println("Server is listening on :8080")
	server.Configure()
	if err := server.Run(":8080"); err != nil {
		log.Fatal("Error starting server")
	}
}
