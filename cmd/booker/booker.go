package main

import (
	"github.com/i-hate-nicknames/redeamtask/pkg/api"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
	"github.com/i-hate-nicknames/redeamtask/pkg/webservice"
)

func main() {
	bookDB := db.MakeMemoryDB()
	API := api.NewAPI(bookDB)
	service := webservice.MakeService(API)
	service.NewBook()
}
