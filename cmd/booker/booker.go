package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"github.com/i-hate-nicknames/redeamtask/pkg/api"
	"github.com/i-hate-nicknames/redeamtask/pkg/db"
	"github.com/i-hate-nicknames/redeamtask/pkg/webservice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

const postgresPort = 5432

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	dbLogger := log.With().Str("component", "db").Logger()
	serviceLogger := log.With().Str("component", "webservice").Logger()
	APILogger := log.With().Str("component", "api").Logger()
	levelVal := os.Getenv("LOG_LEVEL")
	level, err := zerolog.ParseLevel(levelVal)
	if levelVal != "" && err == nil {
		log.Info().Msgf("Log global level is set to: %s", level)
		zerolog.SetGlobalLevel(level)
	} else {
		log.Info().Msgf("Log level is not set, falling back to default: %s", zerolog.GlobalLevel())
	}

	ctx := context.Background()
	bookDB := initDB(ctx, dbLogger)
	defer func() {
		err := bookDB.Close(context.Background())
		if err != nil {
			log.Printf("Error while closing db connection: %s", err)
		}
	}()
	err = bookDB.Migrate(ctx)
	if err != nil {
		log.Fatal().Err(err).Send()
		return
	}
	bookAPI := api.NewAPI(bookDB, APILogger)
	service := webservice.MakeService(bookAPI, serviceLogger)
	port := readEnv("APP_PORT")
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), service)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

// Choose database implementation based on env variable
// and initialize it
func initDB(ctx context.Context, logger zerolog.Logger) db.BookDB {
	dbImpl := os.Getenv("DB")
	var bookDB db.BookDB
	if dbImpl == "memory" {
		bookDB = db.MakeMemoryDB(logger)
	} else {
		var err error
		bookDB, err = db.MakePostgresDB(ctx, dsnFromEnv(), logger)
		if err != nil {
			log.Fatal().Stack().Err(err).Send()
		}
	}
	return bookDB
}

// Read dsn values from env variables or exit program with failure exit code
func dsnFromEnv() string {
	user := readEnv("POSTGRES_USER")
	password := readEnv("POSTGRES_PASSWORD")
	dbName := readEnv("POSTGRES_DB")
	dsnTemplate := "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
	return fmt.Sprintf(dsnTemplate, "db", postgresPort, user, password, dbName)
}

// Read an env variables or exit program with failure exit code
func readEnv(varName string) string {
	val := os.Getenv(varName)
	if val == "" {
		msg := fmt.Sprintf("Missing %s env variable", varName)
		log.Fatal().Stack().Err(errors.New(msg)).Send()
	}
	return val
}
