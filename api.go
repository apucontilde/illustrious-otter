package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/apucontilde/illustrious-otter/database"
)

func validateTransactionRequestBody[T any](body io.Reader, transaction *T) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&transaction)
	if err != nil {
		return err
	}
	return nil
}

func Serve(port string, db *sql.DB) {
	mux := &http.ServeMux{}

	repository := database.NewSQLiteRepository(db)
	mux.Handle("/transaction/{id}", TransactionHandler{transactions: repository})
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
