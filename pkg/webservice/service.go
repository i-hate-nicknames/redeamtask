package webservice

import bookapi "github.com/i-hate-nicknames/redeamtask/pkg/api"

type Service interface {
	Run(port uint16) error
}

type webservice struct {
	api bookapi.BookAPI
}

func MakeService(api bookapi.BookAPI) Service {
	return &webservice{api: api}
}

func (ws *webservice) Run(port uint16) error {
	panic("not implemented")
}
