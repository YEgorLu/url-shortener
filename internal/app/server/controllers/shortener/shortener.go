package shortener

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/storage"
)

type ShortenerController struct {
	group *http.ServeMux
}

func NewShortenerController() *ShortenerController {
	return &ShortenerController{http.NewServeMux()}
}

func (c ShortenerController) RegisterRoute(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == http.MethodPost {
			c.Register(w, r)
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
	newCode, err := storage.Instance().AddURL(url)
	fmt.Println("code ", newCode)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + newCode))
}

func (c ShortenerController) Redirect(w http.ResponseWriter, r *http.Request) {
	if len(r.URL.Path) <= 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := r.URL.Path[1:]
	url, err := storage.Instance().GetURLByCode(code)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
