package transaction

import (
	"context"
	"database/sql"
	"errors"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS transactions(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		userid INTEGER NOT NULL,
		orderid INTEGER NOT NULL,
		storeid INTEGER NOT NULL,
		amount INTEGER NOT NULL,
		details TEXT NOT NULL,
		at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Create(transaction TransactionCreate) (*Transaction, error) {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	res, err := tx.Exec("INSERT INTO transactions( userid, orderid, storeid, amount, details) values(?,?,?,?,?)",
		transaction.UserId, transaction.OrderId, transaction.StoreId, transaction.Amount, transaction.Details)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	row := tx.QueryRow("SELECT * FROM transactions WHERE Id = ?", id)
	var t Transaction
	if err := row.Scan(&t.ID, &t.UserId, &t.OrderId, &t.StoreId, &t.Amount, &t.Details, &t.At); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *SQLiteRepository) All() ([]Transaction, error) {
	rows, err := r.db.Query("SELECT * FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.UserId, &t.OrderId, &t.StoreId, &t.Amount, &t.Details, &t.At); err != nil {
			return nil, err
		}
		all = append(all, t)
	}
	return all, nil
}

func (r *SQLiteRepository) GetById(Id string) (*Transaction, error) {
	row := r.db.QueryRow("SELECT * FROM transactions WHERE Id = ?", Id)

	var t Transaction
	if err := row.Scan(&t.ID, &t.UserId, &t.OrderId, &t.StoreId, &t.Amount, &t.Details, &t.At); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &t, nil
}

func (r *SQLiteRepository) Update(id int64, updated Transaction) (*Transaction, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}
	res, err := r.db.Exec("UPDATE transactions SET userid = ?, orderid = ?, storeid = ?, amount = ?, details = ?  WHERE id = ?",
		&updated.UserId, &updated.OrderId, &updated.StoreId, &updated.Amount, &updated.Details, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

func (r *SQLiteRepository) Delete(id int64) error {
	res, err := r.db.Exec("DELETE FROM transactions WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
