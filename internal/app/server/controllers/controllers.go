package controllers

import (
	"net/http"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/server/controllers/shortener"
)

type controller interface {
	RegisterRoute(mux *http.ServeMux)
}

var controllers = []controller{
	shortener.NewShortenerController(),
}

func Register(mux *http.ServeMux) {
	for _, controller := range controllers {
		controller.RegisterRoute(mux)
	}
}
