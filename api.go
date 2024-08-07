package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/apucontilde/illustrious-otter/transaction"
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
	mux.HandleFunc("/transaction/{id}", HandleTransactionById(repository))
	mux.HandleFunc("/transaction", HandlePOSTTransaction(repository))
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
