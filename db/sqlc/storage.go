package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Storage interface {
	Querier
	ExecuteWithTx(ctx context.Context, f func(*Queries) error) error
}

type SQLStorage struct {
	Querier
	db *sql.DB
}

func NewSQLStorage(db *sql.DB) *SQLStorage {
	return &SQLStorage{Querier: New(db), db: db}
}
func (storage *SQLStorage) ExecuteWithTx(ctx context.Context, f func(*Queries) error) error {
	tx, err := storage.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	txErr := f(q)
	if txErr != nil {
		err := tx.Rollback()
		if err != nil {
			return fmt.Errorf("transaction err: %v, rollback err: %v", txErr, err)
		}
	}

	return tx.Commit()
}
