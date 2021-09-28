package webservice

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/i-hate-nicknames/redeamtask/pkg/api"
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
