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
	r := httptest.NewRequest(http.MethodPost, "/book", body)

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
	publishDate := time.Now()
	b := book.Book{
		Title:       "title",
		Author:      "author",
		Publisher:   "publisher",
		PublishDate: publishDate,
		Rating:      1,
		Status:      book.StatusCheckedOut,
	}
	payload, _ := json.Marshal(b)
	body := bytes.NewBuffer(payload)
	r := httptest.NewRequest(http.MethodPost, "/book", body)

	service.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, w.Body.String())

	var resp IdResponse
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

// workaround for losing precision when marshalling/unmarshalling
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
	return bytes.Compare(t1, t2) == 0
}
