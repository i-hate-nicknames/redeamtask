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

type webservice struct {
	http.Handler
	api bookapi.BookAPI
}

// IDResponse represents json response that includes bookID
type IDResponse struct {
	BookID int `json:"book_id"`
}

// MakeService creates a new web service around book API
func MakeService(api bookapi.BookAPI) http.Handler {
	service := &webservice{api: api}
	r := chi.NewRouter()
	defineRoutes(service, r)
	service.Handler = r
	return service
}

func defineRoutes(ws *webservice, r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/book", func(r chi.Router) {
		r.Post("/", ws.CreateBook)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ws.GetBook)
			r.Delete("/", ws.DeleteBook)
			r.Put("/", ws.UpdateBook)
		})
	})
}

// CreateBook handler creates a new book
func (ws *webservice) CreateBook(w http.ResponseWriter, r *http.Request) {
	b, ok := readBookRequest(w, r)
	if !ok {
		return
	}
	ID, err := ws.api.Create(r.Context(), b)
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := IDResponse{BookID: ID}
	respondJSON(w, http.StatusOK, response)
}

// CreateBook handler updates an existing book
func (ws *webservice) UpdateBook(w http.ResponseWriter, r *http.Request) {
	b, ok := readBookRequest(w, r)
	if !ok {
		return
	}
	IDStr := chi.URLParam(r, "id")
	ID, err := strconv.Atoi(IDStr)
	if err != nil || ID <= 0 {
		respondError(w, http.StatusBadRequest, "book id should be a positive integer")
		return
	}
	err = ws.api.Update(r.Context(), ID, b)
	if err == db.ErrBookNotFound {
		respondError(w, http.StatusNotFound, "book not found")
		return
	}
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// readBookRequest reads book from the given request. In case of error returns
// false and sends error to the client
func readBookRequest(w http.ResponseWriter, r *http.Request) (*book.Book, bool) {
	data, err := ioutil.ReadAll(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Printf("Failed to close body: %s", err)
		}
	}()
	if err != nil {
		log.Printf("Failed to read request: %s", err)
		respondError(w, http.StatusInternalServerError, "backend error")
		return nil, false
	}
	var b book.Book
	err = json.Unmarshal(data, &b)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid book data format")
		return nil, false
	}
	err = b.Validate()
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return nil, false
	}
	return &b, true
}

// GetBook handler looks up a book and sends it to the client
func (ws *webservice) GetBook(w http.ResponseWriter, r *http.Request) {
	IDStr := chi.URLParam(r, "id")
	ID, err := strconv.Atoi(IDStr)
	if err != nil || ID <= 0 {
		respondError(w, http.StatusBadRequest, "book id should be a positive integer")
		return
	}
	book, err := ws.api.Get(r.Context(), ID)
	if err == db.ErrBookNotFound {
		respondError(w, http.StatusNotFound, "book not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, book)
}

// DeleteBook handler deletes a book
func (ws *webservice) DeleteBook(w http.ResponseWriter, r *http.Request) {
	IDStr := chi.URLParam(r, "id")
	ID, err := strconv.Atoi(IDStr)
	if err != nil || ID <= 0 {
		respondError(w, http.StatusBadRequest, "book id should be a positive integer")
		return
	}
	err = ws.api.Delete(r.Context(), ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error string `json:"error"`
}

func respondError(w http.ResponseWriter, status int, msg string) {
	response := ErrorResponse{Error: msg}
	respondJSON(w, status, response)
}

func respondJSON(w http.ResponseWriter, status int, val interface{}) {
	payload, err := json.Marshal(&val)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err != nil {
		panic("failed to marshal error")
	}
	_, err = w.Write(payload)
	if err != nil {
		log.Printf("Error when writing response: %s", err)
	}
}
