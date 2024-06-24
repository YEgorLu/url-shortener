package shortener

import (
	"encoding/json"
	"net/http"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/config"
	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/models/shorten"
	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/server/middleware"
	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/storage"
	"github.com/labstack/echo/v4"
)

type ShortenerController struct {
	stor storage.Storage
}

func NewShortenerController() *ShortenerController {
	return &ShortenerController{
		storage.Instance(),
	}
}

func (c ShortenerController) RegisterRoute(mux *echo.Group) {
	mux.POST("/shorten", c.Register, middleware.UseContentType("application/json"))
	mux.GET("/:code", c.Redirect)
}

func (_c ShortenerController) Register(c echo.Context) error {
	r := c.Request()
	w := c.Response().Writer
	var body shorten.Request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return echo.ErrBadRequest
	}
	if body.Url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return echo.ErrBadRequest
	}
	newCode, err := _c.stor.AddURL(body.Url)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return echo.ErrBadRequest
	}
	c.Response().Header().Set("Content-Type", "application/json")
	response := shorten.Response{config.Params.ShortUrlPrefix + "/" + newCode}
	return c.JSON(http.StatusCreated, response)
}

func (_c ShortenerController) Redirect(c echo.Context) error {
	w := c.Response().Writer
	code := c.Param("code")
	if len(code) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return echo.ErrBadRequest
	}
	url, err := _c.stor.GetURLByCode(code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return echo.ErrBadRequest
	}
	return c.Redirect(http.StatusTemporaryRedirect, url)
}
