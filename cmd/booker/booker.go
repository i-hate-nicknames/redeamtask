package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/i-hate-nicknames/redeamtask/pkg/api"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
	"github.com/i-hate-nicknames/redeamtask/pkg/webservice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	// todo: different logger depending on execution context
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	levelVal := os.Getenv("LOG_LEVEL")
	level, err := zerolog.ParseLevel(levelVal)
	if levelVal != "" && err == nil {
		log.Info().Msgf("Log global level is set to: %s", level)
		zerolog.SetGlobalLevel(level)
	} else {
		log.Info().Msgf("Log level is not set, falling back to default: %s", zerolog.GlobalLevel())
	}
	bookDB, err := db.MakePostgresDB(ctx, dsnFromEnv())
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}
	defer bookDB.Close(context.Background())
	err = bookDB.Migrate(ctx)
	if err != nil {
		log.Fatal().Err(err).Send()
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
		log.Fatal().Err(fmt.Errorf("Missing %s env variable", varName)).Send()
	}
	return val
}
