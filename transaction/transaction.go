package transaction

type TransactionCreate struct {
	Title       string
	Description string
	Ammount     int64
	TxId        string
}

type Transaction struct {
	ID          int64
	Title       string
	Description string
	Ammount     int64
	TxId        string
}
