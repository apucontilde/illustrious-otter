package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/apucontilde/illustrious-otter/transaction"
)

func HandleTransactionById(repository *transaction.SQLiteRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}
func HandlePOSTTransaction(repository *transaction.SQLiteRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

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
			// } else if r.Method == "GET" {
			// 	transactions, err := repository.All()
			// 	if err != nil {
			// 		Error(context.Background(), err.Error())
			// 		return
			// 	}
			// 	Debug("got transactions: %+v\n", transactions)
			// 	w.WriteHeader(http.StatusOK)
			// 	w.Header().Set("Content-Type", "application/json")
			// 	json.NewEncoder(w).Encode(transactions)
			// }
		}
	}
}
