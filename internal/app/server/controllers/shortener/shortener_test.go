package shortener

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/storage"
	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/util"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var _ = storage.Storage(&testStorage{})

type testStorage struct {
	lastURL  string
	lastCode string
}

// AddURL implements storage.Storage.
func (s *testStorage) AddURL(url string) (string, error) {
	if s.lastURL == url {
		return "", errors.New("url exists")
	}
	s.lastURL = url
	s.lastCode = util.RandStringRunes(10)
	return s.lastCode, nil
}

// GetURLByCode implements storage.Storage.
func (s *testStorage) GetURLByCode(code string) (string, error) {
	if s.lastCode == code {
		return s.lastURL, nil
	}
	return "", errors.New("code exists")
}

func (s testStorage) Add() {}

func TestShortenerController_Register(t *testing.T) {
	type args struct {
		url string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "initial url", args: args{"http://some.url.com"}, want: want{http.StatusCreated}},
		{name: "existing url", args: args{"http://some.url.com"}, want: want{http.StatusBadRequest}},
		{name: "new url", args: args{"hhtp://some.other.url.com"}, want: want{http.StatusCreated}},
	}

	c := NewShortenerController()
	c.stor = &testStorage{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.args.url))
			e := echo.New()

			ctx := e.NewContext(r, w)
			defer e.ReleaseContext(ctx)

			c.Register(ctx)

			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)

			t.Log("body ", w.Body.String(), res.StatusCode)

			if w.Code != http.StatusCreated {
				var body []byte
				if _, err := res.Body.Read(body); err != nil && err != io.EOF {
					t.Fatal("can't read body")
				}
				defer res.Body.Close()
				assert.Empty(t, body)
			} else {
				assert.Equal(t, "text/plain", res.Header.Get("Content-Type"))
				assert.Regexp(t, "/[0-9_a-zA-Z]{10}$", w.Body.String())
			}
		})
	}
}

func TestShortenerController_Redirect(t *testing.T) {
	c := NewShortenerController()
	c.stor = &testStorage{}

	type args struct {
		url      string
		register bool
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{name: "normal use", args: args{"http://some.url.com", true}, want: want{http.StatusTemporaryRedirect}},
		{name: "another use", args: args{"http://some.another.url.com", true}, want: want{http.StatusTemporaryRedirect}},
		{name: "no registered code", args: args{"", false}, want: want{http.StatusBadRequest}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := "somerandomcode"
			e := echo.New()
			if tt.args.register {
				wReg := httptest.NewRecorder()
				rReg := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.args.url))
				ctx1 := e.NewContext(rReg, wReg)
				defer e.ReleaseContext(ctx1)
				c.Register(ctx1)
				redirectUrlBytes, err := io.ReadAll(wReg.Result().Body)
				if err != nil {
					t.Fatal("Error reading response body")
				}
				redirectUrl := string(redirectUrlBytes)
				code = strings.Split(redirectUrl, "/")[3] // http://localhost:8080/TvRWs1XPcW
			}
			wRed := httptest.NewRecorder()
			rRed := httptest.NewRequest(http.MethodGet, "/"+code, nil)
			ctx2 := e.NewContext(rRed, wRed)
			ctx2.SetParamNames("code")
			ctx2.SetParamValues(code)
			defer e.ReleaseContext(ctx2)
			c.Redirect(ctx2)

			assert.Equal(t, tt.want.code, wRed.Code)
		})
	}
}
