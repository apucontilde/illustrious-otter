package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"errors"
	"io"
	"net/http"

	"github.com/apucontilde/illustrious-otter/transaction"
	_ "github.com/mattn/go-sqlite3"
)

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

	mux.HandleFunc("/transaction/:id", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		transaction, err := repository.GetById(r.PathValue("id"))
		if err != nil {
			panic(err)
		}
		fmt.Printf("got transaction: %+v\n", transaction)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transaction)

	})
	mux.HandleFunc("/transaction", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var transactionCreate transaction.TransactionCreate
			err := validateTransactionRequestBody(r.Body, &transactionCreate)
			if err != nil {
				fmt.Printf("%s\n", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			transaction, err := repository.Create(transactionCreate)
			if err != nil {
				fmt.Printf("%s\n", err)
				if err.Error() == "record already exists" {
					w.WriteHeader(http.StatusForbidden)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				json.NewEncoder(w).Encode(struct {
					ErrorMessage string
				}{
					ErrorMessage: err.Error(),
				})
				return
			}
			fmt.Printf("created transaction: %+v\n", transaction)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(transaction)
		} else if r.Method == "GET" {
			transactions, err := repository.All()
			if err != nil {
				panic(err)
			}
			fmt.Printf("got transaction: %+v\n", transactions)
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(transactions)
		}
	})
	server := &http.Server{
		Addr:    ":3333",
		Handler: mux,
	}
	fmt.Printf("listening on %s\n", server.Addr)
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error listening for server: %s\n", err)
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
			log.Fatal(err)
		}

	}

	Serve(":3333", transactionRepository)

}
