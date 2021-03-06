package webservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/i-hate-nicknames/redeamtask/pkg/api"
	"github.com/i-hate-nicknames/redeamtask/pkg/book"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
)

// create new service with given books, return service
// and slice of book ids, where n-th element corresponds to
// assigned book ID to the n-th book in the provided books slice
func newService(books ...db.BookRecord) (http.Handler, []int) {
	dbLogger := log.With().Str("component", "test_db").Logger()
	serviceLogger := log.With().Str("component", "test_webservice").Logger()
	APILogger := log.With().Str("component", "test_api").Logger()
	memdb := db.MakeMemoryDB(dbLogger)
	var ids []int
	for _, book := range books {
		record, _ := memdb.Create(context.Background(), book) // nolint: errcheck
		ids = append(ids, record.ID)
	}
	api := api.NewAPI(memdb, APILogger)
	return MakeService(api, serviceLogger), ids
}

func TestInvalidPath(t *testing.T) {
	service, _ := newService()
	w := httptest.NewRecorder()
	body := bytes.NewBuffer(make([]byte, 0))
	r := httptest.NewRequest(http.MethodGet, "/INVALID", body)

	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
}

func TestCreateBookCorrectID(t *testing.T) {
	service, _ := newService()
	w := httptest.NewRecorder()
	b := book.Book{
		Title:       "title",
		Author:      "author",
		Publisher:   "publisher",
		PublishDate: time.Now(),
		Rating:      1,
		Status:      book.StatusCheckedOut,
	}
	payload, _ := json.Marshal(b) // nolint: errcheck
	body := bytes.NewBuffer(payload)
	r := httptest.NewRequest(http.MethodPost, "/book", body)

	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp IDResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.BookID)
}

func TestCreateBookCorrectData(t *testing.T) {
	service, _ := newService()
	w := httptest.NewRecorder()
	publishDate := time.Now()
	b := book.Book{
		Title:       "title",
		Author:      "author",
		Publisher:   "publisher",
		PublishDate: publishDate,
		Rating:      1,
		Status:      book.StatusCheckedOut,
	}
	payload, _ := json.Marshal(b) // nolint: errcheck
	body := bytes.NewBuffer(payload)
	r := httptest.NewRequest(http.MethodPost, "/book", body)

	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp IDResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	body = bytes.NewBuffer(make([]byte, 0))
	path := fmt.Sprintf("/book/%d", resp.BookID)
	r = httptest.NewRequest(http.MethodGet, path, body)
	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var bookResp book.Book
	err = json.Unmarshal(w.Body.Bytes(), &bookResp)
	assert.NoError(t, err)
	assert.True(t, compareBooks(b, bookResp))
}

func TestGetBook(t *testing.T) {
	publishDate := time.Now()
	b := book.Book{
		Title:       "title",
		Author:      "author",
		Publisher:   "publisher",
		PublishDate: publishDate,
		Rating:      1,
		Status:      book.StatusCheckedOut,
	}
	br := db.BookRecord{Book: b}
	service, ids := newService(br)

	w := httptest.NewRecorder()
	body := bytes.NewBuffer(make([]byte, 0))
	path := fmt.Sprintf("/book/%d", ids[0])
	r := httptest.NewRequest(http.MethodGet, path, body)
	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var bookResp book.Book
	err := json.Unmarshal(w.Body.Bytes(), &bookResp)
	assert.NoError(t, err)
	assert.True(t, compareBooks(b, bookResp))
}

func TestDeleteBook(t *testing.T) {
	publishDate := time.Now()
	b := book.Book{
		Title:       "title",
		Author:      "author",
		Publisher:   "publisher",
		PublishDate: publishDate,
		Rating:      1,
		Status:      book.StatusCheckedOut,
	}
	br := db.BookRecord{Book: b}
	service, ids := newService(br)

	w := httptest.NewRecorder()
	body := bytes.NewBuffer(make([]byte, 0))
	path := fmt.Sprintf("/book/%d", ids[0])
	r := httptest.NewRequest(http.MethodDelete, path, body)
	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	r = httptest.NewRequest(http.MethodGet, path, body)
	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
}

func TestUpdateBook(t *testing.T) {
	publishDate := time.Now()
	b := book.Book{
		Title:       "title",
		Author:      "author",
		Publisher:   "publisher",
		PublishDate: publishDate,
		Rating:      1,
		Status:      book.StatusCheckedOut,
	}
	br := db.BookRecord{Book: b}
	service, ids := newService(br)

	b2 := book.Book{
		Title:       "title2",
		Author:      "author2",
		Publisher:   "publisher2",
		PublishDate: publishDate,
		Rating:      2,
		Status:      book.StatusCheckedOut,
	}
	w := httptest.NewRecorder()
	payload, _ := json.Marshal(b2) // nolint: errcheck
	body := bytes.NewBuffer(payload)
	path := fmt.Sprintf("/book/%d", ids[0])
	r := httptest.NewRequest(http.MethodPut, path, body)
	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	r = httptest.NewRequest(http.MethodGet, path, body)
	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var bookResp book.Book
	err := json.Unmarshal(w.Body.Bytes(), &bookResp)
	assert.NoError(t, err)
	assert.True(t, compareBooks(b2, bookResp))
}

// workaround for losing precision when marshaling/unmarshalling
// time.Time values.
// todo: investigate and fix
func compareBooks(b1, b2 book.Book) bool {
	if b1.Title != b2.Title {
		return false
	}
	if b1.Author != b2.Author {
		return false
	}
	if b1.Rating != b2.Rating {
		return false
	}
	if b1.Status != b2.Status {
		return false
	}
	if b1.Publisher != b2.Publisher {
		return false
	}
	t1, err := b1.PublishDate.MarshalText()
	if err != nil {
		panic("marshaling error")
	}
	t2, err := b2.PublishDate.MarshalText()
	if err != nil {
		panic("marshaling error")
	}
	return bytes.Equal(t1, t2)
}
