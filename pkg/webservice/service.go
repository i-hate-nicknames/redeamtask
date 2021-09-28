package webservice

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	bookapi "github.com/i-hate-nicknames/redeamtask/pkg/api"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
)

type Service interface {
	http.Handler
}

type webservice struct {
	http.Handler
	api bookapi.BookAPI
}

type IdResponse struct {
	BookID int `json:"book_id"`
}

func MakeService(api bookapi.BookAPI) Service {
	service := &webservice{api: api}
	r := chi.NewRouter()
	defineRoutes(service, r)
	service.Handler = r
	return service
}

func defineRoutes(ws *webservice, r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Route("/books/{id}", func(r chi.Router) {
		r.Get("/", ws.GetBook)
	})
}

func (ws *webservice) GetBook(w http.ResponseWriter, r *http.Request) {
	IDStr := chi.URLParam(r, "id")
	ID, err := strconv.Atoi(IDStr)
	if err != nil || ID <= 0 {
		respondError(w, http.StatusBadRequest, "book id should be a positive integer")
		return
	}
	book, err := ws.api.Get(ID)
	if err == db.ErrBookNotFound {
		respondError(w, http.StatusNotFound, "book not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJson(w, http.StatusOK, book)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func respondError(w http.ResponseWriter, status int, msg string) {
	response := ErrorResponse{Error: msg}
	respondJson(w, status, response)
}

func respondJson(w http.ResponseWriter, status int, val interface{}) {
	payload, err := json.Marshal(&val)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err != nil {
		panic("failed to marshal error")
	}
	w.Write(payload)
}
