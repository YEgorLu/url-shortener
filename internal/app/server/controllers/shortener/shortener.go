package shortener

import (
	"fmt"
	"io"
	"net/http"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/server/middleware"
	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
	mux.POST("/", c.Register, middleware.UseContentType("text/html"))
	mux.GET("/:code", c.Redirect)
}

func (_c ShortenerController) Register(c echo.Context) error {
	r := c.Request()
	w := c.Response().Writer
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return echo.ErrBadRequest
	}
	url := string(body)
	fmt.Println("url ", url)
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return echo.ErrBadRequest
	}
	newCode, err := _c.stor.AddURL(url)
	fmt.Println("code ", newCode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return echo.ErrBadRequest
	}
	c.Response().Header().Set("Content-Type", "text/plain")
	return c.String(http.StatusCreated, "http://localhost:8080/"+newCode)
	// w.Header().Set("Content-Type", "text/plain")
	// w.WriteHeader(http.StatusCreated)
	// w.Write([]byte("http://localhost:8080/" + newCode))
	//return nil
}

func (_c ShortenerController) Redirect(c echo.Context) error {
	//r := c.Request()
	w := c.Response().Writer
	code := c.Param("code")
	log.Info("code ", code)
	if len(code) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return echo.ErrBadRequest
	}
	url, err := _c.stor.GetURLByCode(code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return echo.ErrBadRequest
	}
	//http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return c.Redirect(http.StatusTemporaryRedirect, url)
}
