package shortener

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/server/middleware"
	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/storage"
)

type ShortenerController struct {
	group *http.ServeMux
	stor  storage.Storage
}

func NewShortenerController() *ShortenerController {
	return &ShortenerController{
		http.NewServeMux(),
		storage.Instance(),
	}
}

func (c ShortenerController) RegisterRoute(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == http.MethodPost {
			middleware.Use(http.HandlerFunc(c.Register), middleware.UseContentType("text/html")).ServeHTTP(w, r)
		} else if r.Method == http.MethodGet && len(strings.SplitN(r.URL.Path, "/", 3)) == 2 {
			c.Redirect(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func (c ShortenerController) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url := string(body)
	fmt.Println("url ", url)
	if url == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newCode, err := c.stor.AddURL(url)
	fmt.Println("code ", newCode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + newCode))
}

func (c ShortenerController) Redirect(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) <= 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := r.URL.Path[1:]
	url, err := c.stor.GetURLByCode(code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
