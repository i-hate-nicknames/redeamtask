package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/i-hate-nicknames/redeamtask/pkg/api"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
	"github.com/i-hate-nicknames/redeamtask/pkg/webservice"
)

func main() {
	bookDB, err := db.MakePostgresDB(context.Background(), dsnFromEnv())
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

// Read dsn values from env variables or exit program with failure exit code
func dsnFromEnv() string {
	port := readEnv("POSTGRES_PORT")
	user := readEnv("POSTGRES_USER")
	password := readEnv("POSTGRES_PASSWORD")
	dbName := readEnv("POSTGRES_DB")
	dsnTemplate := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	return fmt.Sprintf(dsnTemplate, "db", port, user, password, dbName)
}

// Read an env variables or exit program with failure exit code
func readEnv(varName string) string {
	val := os.Getenv(varName)
	if val == "" {
		log.Fatalf("Missing %s env variable", varName)
	}
	return val
}
