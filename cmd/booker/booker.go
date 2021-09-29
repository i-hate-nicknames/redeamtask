package main

import (
	"net/http"

	"github.com/i-hate-nicknames/redeamtask/pkg/api"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
	"github.com/i-hate-nicknames/redeamtask/pkg/webservice"
)

func main() {
	bookDB := db.MakeMemoryDB()
	defer bookDB.Close()
	bookAPI := api.NewAPI(bookDB)
	service := webservice.MakeService(bookAPI)
	http.ListenAndServe(":8080", service)
}
