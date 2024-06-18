package storage

import "github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/storage/localStorage"

type Storage interface {
	AddURL(url string) (string, error)
	GetURLByCode(code string) (string, error)
}

var instance Storage

func Instance() Storage {
	if instance == nil {
		instance = localStorage.NewLocalStorage()
	}
	return instance
}
