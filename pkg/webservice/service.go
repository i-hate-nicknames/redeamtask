package webservice

import (
	"net/http"

	bookapi "github.com/i-hate-nicknames/redeamtask/pkg/api"
)

type Service interface {
	http.Handler
	NewBook()
}

type webservice struct {
	api bookapi.BookAPI
}

func MakeService(api bookapi.BookAPI) Service {
	panic("not implemented")
}

func (ws *webservice) Run(port uint16) error {
	panic("not implemented")
}
