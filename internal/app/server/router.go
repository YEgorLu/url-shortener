package server

import (
	"net/http"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/server/controllers"
	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/server/middleware"
)

func newRouter() http.Handler {
	router := http.NewServeMux()
	controllers.Register(router)
	return middleware.Use(router, middleware.UseLogging())
}
