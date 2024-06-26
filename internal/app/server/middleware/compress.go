package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
)

var _ = http.ResponseWriter(&compressWriter{})

type compressWriter struct {
	w  http.ResponseWriter
	gz *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		gz: gzip.NewWriter(w),
	}
}

func (w *compressWriter) Header() http.Header {
	return w.w.Header()
}

func (w *compressWriter) WriteHeader(code int) {
	w.w.WriteHeader(code)
}

func (w *compressWriter) Write(data []byte) (int, error) {
	return w.gz.Write(data)
}

var _ = io.ReadCloser(&compressBody{})

type compressBody struct {
	r  io.ReadCloser
	gz *gzip.Reader
}

func newCompressBody(r io.ReadCloser) (*compressBody, error) {
	reader, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return &compressBody{
		r:  r,
		gz: reader,
	}, nil
}

func (b *compressBody) Close() error {
	if err := b.r.Close(); err != nil {
		return err
	}
	return b.gz.Close()
}

func (b *compressBody) Read(p []byte) (n int, err error) {
	return b.gz.Read(p)
}

func Compression() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return echo.HandlerFunc(func(c echo.Context) error {
			h := c.Request().Header
			w := c.Response().Writer
			if slices.Contains(h.Values(echo.HeaderAcceptEncoding), "gzip") && compressContentType(h.Get(echo.HeaderContentType)) {
				wr := gzip.NewWriter(w)
				defer wr.Close()
				c.Response().Writer = &compressWriter{w, wr}
			}
			if slices.Contains(h.Values(echo.HeaderContentEncoding), "gzip") {
				body := c.Request().Body
				cb, err := newCompressBody(body)
				if err != nil {
					return echo.ErrUnsupportedMediaType
				}
				defer cb.Close()
				c.Request().Body = cb
			}
			return next(c)
		})
	}
}

var compressedContentTypes = map[string]bool{
	echo.MIMEApplicationJSON: true,
	echo.MIMETextHTML:        true,
}

func compressContentType(contentType string) bool {
	_, ok := compressedContentTypes[contentType]
	return ok
}
