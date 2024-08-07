package main

import (
	"database/sql"
	"log"
	"os"

	"context"
	"errors"

	"log/slog"

	"github.com/apucontilde/illustrious-otter/transaction"
	_ "github.com/mattn/go-sqlite3"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
var Info = logger.Info
var Debug = logger.Debug
var Error = logger.ErrorContext

const db_file = "sqlite.db"

func main() {
	var do_migrate = true
	if _, err := os.Stat(db_file); err == nil {
		do_migrate = false
	} else if errors.Is(err, os.ErrNotExist) {
	} else {
		log.Fatal(err)
	}

	db, err := sql.Open("sqlite3", db_file)
	if err != nil {
		log.Fatal(err)
	}

	transactionRepository := transaction.NewSQLiteRepository(db)

	if do_migrate {
		if err := transactionRepository.Migrate(); err != nil {
			Error(context.Background(), "error running migrations %s\n", err)
		}

	}

	Serve(":3333", transactionRepository)

}
