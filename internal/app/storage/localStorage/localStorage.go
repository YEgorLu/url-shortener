package localStorage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/YEgorLu/go-musthave-shortener-tpl/internal/app/util"
)

type localStorage struct {
	urlToCode *sync.Map
	codeToUrl *sync.Map
}

func NewLocalStorage() *localStorage {
	return &localStorage{
		urlToCode: &sync.Map{},
		codeToUrl: &sync.Map{},
	}
}

func (s localStorage) AddURL(url string) (string, error) {
	if _, ok := s.urlToCode.Load(url); ok {
		return "", errors.New(fmt.Sprint("URL already exists ", url))
	}
	code := util.RandStringRunes(10)
	s.urlToCode.Store(url, code)
	s.codeToUrl.Store(code, url)
	return code, nil
}

func (s localStorage) GetURLByCode(code string) (string, error) {
	if url, ok := s.codeToUrl.Load(code); !ok {
		return "", errors.New(fmt.Sprint("Code does not exists ", code))
	} else {
		var urlStr string
		if urlStr, ok = url.(string); !ok {
			errorTxt := fmt.Sprint("Could not assert url ", url, " to string")
			return "", errors.New(errorTxt)
		}
		return urlStr, nil
	}
}
