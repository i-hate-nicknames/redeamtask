package webservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/i-hate-nicknames/redeamtask/pkg/api"
	"github.com/i-hate-nicknames/redeamtask/pkg/book"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
)

func newService() Service {
	memdb := db.MakeMemoryDB()
	api := api.NewAPI(memdb)
	return MakeService(api)
}

func TestInvalidPath(t *testing.T) {
	service := newService()
	w := httptest.NewRecorder()
	body := bytes.NewBuffer(make([]byte, 0))
	r := httptest.NewRequest(http.MethodGet, "/INVALID", body)

	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code, w.Body.String())
}

func TestCreateBookCorrectID(t *testing.T) {
	service := newService()
	w := httptest.NewRecorder()
	b := book.Book{
		Title:       "title",
		Author:      "author",
		Publisher:   "publisher",
		PublishDate: time.Now(),
		Rating:      1,
		Status:      book.StatusCheckedOut,
	}
	payload, _ := json.Marshal(b)
	body := bytes.NewBuffer(payload)
	r := httptest.NewRequest(http.MethodPost, "/books", body)

	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp IdResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.BookID)
}

func TestCreateBookCorrectData(t *testing.T) {
	service := newService()
	w := httptest.NewRecorder()
	b := book.Book{
		Title:       "title",
		Author:      "author",
		Publisher:   "publisher",
		PublishDate: time.Now(),
		Rating:      1,
		Status:      book.StatusCheckedOut,
	}
	payload, _ := json.Marshal(b)
	body := bytes.NewBuffer(payload)
	r := httptest.NewRequest(http.MethodPost, "/books", body)

	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp IdResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)

	w = httptest.NewRecorder()
	body = bytes.NewBuffer(make([]byte, 0))
	path := fmt.Sprintf("/books/%d", resp.BookID)
	r = httptest.NewRequest(http.MethodGet, path, body)
	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var bookResp book.Book
	err = json.Unmarshal(w.Body.Bytes(), &bookResp)
	assert.NoError(t, err)
	assert.Equal(t, b, bookResp)
}
