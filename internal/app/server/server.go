package server

import (
	"net/http"
)

var server http.Server

func Configure() {
	server = http.Server{
		Handler: newRouter(),
	}
}

func Run(addr string) error {
	server.Addr = addr
	return server.ListenAndServe()
}
