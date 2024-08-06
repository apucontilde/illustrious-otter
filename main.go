package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"context"
	"errors"
	"io"
	"net/http"

	"log/slog"

	"github.com/apucontilde/illustrious-otter/transaction"
	_ "github.com/mattn/go-sqlite3"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
var Info = logger.Info
var Debug = logger.Debug
var Error = logger.ErrorContext

func validateTransactionRequestBody(body io.Reader, transaction *transaction.TransactionCreate) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&transaction)
	if err != nil {
		return err
	}
	return nil
}

func Serve(port string, repository *transaction.SQLiteRepository) {
	mux := &http.ServeMux{}

	mux.HandleFunc("/transaction/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		transaction, err := repository.GetById(r.PathValue("id"))
		if err != nil {
			errMessage := err.Error()
			if errMessage == "row not exists" {
				w.WriteHeader(http.StatusNotFound)
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			Error(context.Background(), errMessage)

		}
		Debug("get transaction: %+v\n", transaction)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transaction)

	})
	mux.HandleFunc("/transaction", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			var transactionCreate transaction.TransactionCreate
			err := validateTransactionRequestBody(r.Body, &transactionCreate)
			if err != nil {
				Error(context.Background(), err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			Debug("post transaction: %+v\n", transactionCreate)
			transaction, err := repository.Create(transactionCreate)
			if err != nil {
				errMessage := err.Error()
				if errMessage == "record already exists" {
					w.WriteHeader(http.StatusForbidden)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				Error(context.Background(), errMessage)
				return
			}
			Debug("created transaction: %+v\n", transaction)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(transaction)
		} else if r.Method == "GET" {
			transactions, err := repository.All()
			if err != nil {
				Error(context.Background(), err.Error())
				return
			}
			Debug("got transactions: %+v\n", transactions)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(transactions)
		}
	})
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
	}
	Info("listening", "addr", server.Addr)
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		Info("server closed\n")
	} else if err != nil {
		Error(context.Background(), "error listening for server: %s\n", err)
	}
}

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
