package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rberr := tx.Rollback(); rberr != nil {
			return fmt.Errorf("error with transcation rollback: %w", err)
		}
		// TODO wrap the error!
		return fmt.Errorf("%w", err)
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		////////////
		// transfer creating
		////////////
		res, err := q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("cannot create transfer: %w", err)
		}

		tId, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("cannot get last transfer id: %w", err)
		}

		result.Transfer, err = q.GetTransfer(ctx, tId)
		if err != nil {
			return fmt.Errorf("cannot get transfer: %w", err)
		}
		////////////
		//FROM entry creating
		////////////
		res, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("cannot create entry: %w", err)
		}

		eId, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("cannot get last entry id: %w", err)
		}

		result.FromEntry, err = q.GetEntry(ctx, eId)
		if err != nil {
			return fmt.Errorf("cannot get entry: %w", err)
		}
		////////////
		//TO entry creating
		////////////
		res, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("cannot create entry: %w", err)
		}

		eId, err = res.LastInsertId()
		if err != nil {
			return fmt.Errorf("cannot get last entry id: %w", err)
		}

		result.ToEntry, err = q.GetEntry(ctx, eId)
		if err != nil {
			return fmt.Errorf("cannot get entry: %w", err)
		}

		//TODO: update accounts balance

		return nil
	})

	return result, err
}
