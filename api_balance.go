package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/apucontilde/illustrious-otter/database"
)

type BalanceHandler struct {
	transactions *database.SQLiteRepository
}

func (t BalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	storeId := r.PathValue("StoreId")
	if storeId == "" {
		w.WriteHeader(http.StatusBadRequest)
		Error(context.Background(), "No StoreId provided for GetStoreBalance")
		return
	}
	var from string = "0"
	if r.URL.Query().Has("From") {
		from = r.URL.Query().Get("From")
	}
	var to string = "CURRENT_TIMESTAMP"
	if r.URL.Query().Has("To") {
		to = r.URL.Query().Get("To")
	}
	switch r.Method {
	case http.MethodGet:
		balances, err := t.transactions.GetStoreBalance(storeId, from, to)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			Error(context.Background(), err.Error())
			return
		}
		Debug("got balances: %+v\n", balances)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(balances)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
