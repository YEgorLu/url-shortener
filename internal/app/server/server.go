package server

import (
	"fmt"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/config"
	"github.com/labstack/echo/v4"
)

var server echo.Echo

func Configure() {
	server = *newRouter()
}

func Run() error {
	for _, route := range server.Routes() {
		fmt.Printf("path: %s; method: %s; name: %s\r\n", route.Path, route.Method, route.Name)
	}
	return server.Start(config.Params.ServerAddress)
}
