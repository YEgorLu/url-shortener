package server

import (
	"github.com/labstack/echo/v4"
)

var server echo.Echo

func Configure() {
	server = *newRouter()
}

func Run(addr string) error {
	for _, route := range server.Routes() {
		server.Logger.Print(route.Method, "\n", route.Path, "\n", route.Name)
	}
	return server.Start(addr)
}
