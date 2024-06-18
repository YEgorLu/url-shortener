package middleware

import (
	"log"
	"net/http"
	"slices"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/util"
)

type Middleware func(next http.Handler) http.Handler

func Use(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func UseMethod(methodNames ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if slices.Index(methodNames, r.Method) == -1 {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UseContentType(contentTypes ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !util.Intersects(contentTypes, r.Header.Values("Content-Type")) {
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UseLogging() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request %s %s %s", r.Method, r.URL, r.Header.Values("Content-Type"))
			next.ServeHTTP(w, r)
			log.Printf("Response %s %s %s %s", r.Method, r.URL, w.Header().Get("Status"), w.Header().Get("Content-Type"))
		})
	}
}
