package transaction

import "time"

type TransactionCreate struct {
	UserId  int64
	OrderId int64
	StoreId int64
	Details string
	Amount  int64
}

/*
ammount is double precission currency
*/
type Transaction struct {
	ID      int64
	UserId  int64
	OrderId int64
	StoreId int64
	Details string
	Amount  int64
	At      time.Time
}
