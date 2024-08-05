package transaction

import (
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
		title TEXT NOT NULL UNIQUE,
		description TEXT NOT NULL,
		ammount INTEGER NOT NULL,
		txid TEXT NOT NULL UNIQUE
	);
	`

	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Create(transaction Transaction) (*Transaction, error) {
	res, err := r.db.Exec("INSERT INTO transactions(title, description, ammount, txid) values(?,?,?, ?)", transaction.Title, transaction.Description, transaction.Ammount, transaction.TxId)
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
	transaction.ID = id

	return &transaction, nil
}

func (r *SQLiteRepository) All() ([]Transaction, error) {
	rows, err := r.db.Query("SELECT * FROM transactions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.ID, &transaction.Title, &transaction.Description, &transaction.Ammount, &transaction.TxId); err != nil {
			return nil, err
		}
		all = append(all, transaction)
	}
	return all, nil
}

func (r *SQLiteRepository) GetByName(name string) (*Transaction, error) {
	row := r.db.QueryRow("SELECT * FROM transactions WHERE title = ?", name)

	var transaction Transaction
	if err := row.Scan(&transaction.ID, &transaction.Title, &transaction.Description, &transaction.Ammount, &transaction.TxId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *SQLiteRepository) Update(id int64, updated Transaction) (*Transaction, error) {
	if id == 0 {
		return nil, errors.New("invalid updated ID")
	}
	res, err := r.db.Exec("UPDATE transactions SET title = ?, description = ?, ammount = ?, txid = ? WHERE id = ?", updated.Title, updated.Description, updated.Ammount, updated.TxId, id)
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
