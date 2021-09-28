package webservice

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	bookapi "github.com/i-hate-nicknames/redeamtask/pkg/api"
	"github.com/i-hate-nicknames/redeamtask/pkg/book"
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
	r.Route("/books", func(r chi.Router) {
		r.Get("/{id}", ws.GetBook)
		r.Post("/", ws.CreateBook)
	})
}

func (ws *webservice) CreateBook(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Printf("Failed to read request: %s", err)
		respondError(w, http.StatusInternalServerError, "backend error")
		return
	}
	var b book.Book
	err = json.Unmarshal(data, &b)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid book data format")
		return
	}
	err = b.Validate()
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	ID, err := ws.api.Create(&b)
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := IdResponse{BookID: ID}
	respondJson(w, http.StatusOK, response)
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
