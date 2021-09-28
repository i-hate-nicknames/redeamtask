package main

import (
	"net/http"

	"github.com/i-hate-nicknames/redeamtask/pkg/api"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
	"github.com/i-hate-nicknames/redeamtask/pkg/webservice"
)

func main() {
	bookDB := db.MakeMemoryDB()
	API := api.NewAPI(bookDB)
	service := webservice.MakeService(API)
	http.ListenAndServe(":8080", service)
}
