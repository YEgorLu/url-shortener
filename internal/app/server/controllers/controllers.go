package controllers

import (
	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/server/controllers/shortener"
	"github.com/labstack/echo/v4"
)

type controller interface {
	RegisterRoute(mux *echo.Group)
}

var controllers = []controller{
	shortener.NewShortenerController(),
}

func Register(mux *echo.Echo) {
	for _, controller := range controllers {
		controller.RegisterRoute(mux.Group("/api"))
	}
}
