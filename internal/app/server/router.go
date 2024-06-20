package server

import (
	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/server/controllers"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func newRouter() *echo.Echo {
	router := echo.New()
	controllers.Register(router)
	router.Use(echomiddleware.Logger())
	return router
}
