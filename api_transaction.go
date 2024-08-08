package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/apucontilde/illustrious-otter/database"
)

type TransactionHandler struct {
	transactions *database.SQLiteRepository
}

func (t TransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		transaction, err := t.transactions.GetById(r.PathValue("id"))
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
	} else if r.Method == "PATCH" {
		// TODO evaluate the need for transactionUpdate type(null fields)
		// var transactionUpdate transaction.Transaction
		var transactionUpdate map[string]string
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&transactionUpdate)
		if err != nil {
			Error(context.Background(), err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ID := r.PathValue("id")

		var newTransaction *database.Transaction
		if details := transactionUpdate["Details"]; details != "" {
			newTransaction, err = t.transactions.UpdateField(ID, "Details", details)
		}
		if orderId := transactionUpdate["OrderId"]; orderId != "" {
			newTransaction, err = t.transactions.UpdateField(ID, "OrderId", orderId)
		}
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
		Debug("updated transaction: %+v\n", newTransaction)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newTransaction)
	} else if r.Method == "POST" {
		var transactionCreate database.TransactionCreate
		err := validateTransactionRequestBody[database.TransactionCreate](r.Body, &transactionCreate)
		if err != nil {
			Error(context.Background(), err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		Debug("post transaction: %+v\n", transactionCreate)
		transaction, err := t.transactions.Create(transactionCreate)
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
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
