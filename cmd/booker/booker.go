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
	ctx := context.Background()
	bookDB := initDB(ctx)
	fmt.Println("hui sasi")
	defer func() {
		err := bookDB.Close(context.Background())
		if err != nil {
			log.Printf("Error while closing db connection: %s", err)
		}
	}()
	err := bookDB.Migrate(ctx)
	if err != nil {
		log.Fatal(err)
	}
	bookAPI := api.NewAPI(bookDB)
	service := webservice.MakeService(bookAPI)
	port := readEnv("APP_PORT")
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), service)
	if err != nil {
		log.Fatal(err)
	}
}

// Choose database implementation based on env variable
// and initialize it
func initDB(ctx context.Context) db.BookDB {
	dbImpl := os.Getenv("DB")
	var bookDB db.BookDB
	if dbImpl == "memory" {
		bookDB = db.MakeMemoryDB()
	} else {
		var err error
		bookDB, err = db.MakePostgresDB(ctx, dsnFromEnv())
		if err != nil {
			log.Fatal(err)
		}
	}
	return bookDB
}

// Read dsn values from env variables or exit program with failure exit code
func dsnFromEnv() string {
	port := readEnv("POSTGRES_PORT")
	user := readEnv("POSTGRES_USER")
	password := readEnv("POSTGRES_PASSWORD")
	dbName := readEnv("POSTGRES_DB")
	dsnTemplate := "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
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
