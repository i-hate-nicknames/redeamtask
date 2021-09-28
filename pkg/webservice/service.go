package webservice

import (
	"net/http"

	"github.com/go-chi/chi"
	bookapi "github.com/i-hate-nicknames/redeamtask/pkg/api"
)

type Service interface {
	http.Handler
}

type webservice struct {
	http.Handler
	api bookapi.BookAPI
}

func MakeService(api bookapi.BookAPI) Service {
	service := &webservice{api: api}
	r := chi.NewRouter()
	defineRoutes(service, r)
	service.Handler = r
	return service
}

func defineRoutes(ws *webservice, router *chi.Mux) {

}
