package middleware

import (
	"log"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/util"
	"github.com/labstack/echo/v4"
)

func UseContentType(contentTypes ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			if !util.Intersects(contentTypes, c.Request().Header.Values("Content-Type")) {
				return echo.ErrUnsupportedMediaType
			}
			next(c)
			return nil
		})
	}
}

func UseLogging() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			r := c.Request()
			w := c.Response().Writer
			log.Printf("Request %s %s %s", r.Method, r.URL, r.Header.Values("Content-Type"))
			next(c)
			log.Printf("Response %s %s %s %s", r.Method, r.URL, w.Header().Get("Status"), w.Header().Get("Content-Type"))
			return nil
		})
	}
}
