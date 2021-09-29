package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/i-hate-nicknames/redeamtask/pkg/api"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
	"github.com/i-hate-nicknames/redeamtask/pkg/webservice"
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 8035, "root", "pass", "booker")
	bookDB, err := db.MakePostgresDB(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer bookDB.Close()
	err = bookDB.Migrate()
	if err != nil {
		log.Fatal(err)
	}
	bookAPI := api.NewAPI(bookDB)
	service := webservice.MakeService(bookAPI)
	http.ListenAndServe(":8080", service)
}
