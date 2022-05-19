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

var txKey = struct{}{}

func (s *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)

		////////////
		// transfer creating
		////////////
		fmt.Println(txName, "create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("cannot create transfer: %w", err)
		}

		////////////
		//FROM entry creating
		////////////
		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("cannot create entry: %w", err)
		}

		////////////
		//TO entry creating
		////////////
		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("cannot create entry: %w", err)
		}

		//TODO: update accounts balance
		fmt.Println(txName, "get account 1")
		acc1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		if err != nil {
			return fmt.Errorf("cannot get account: %w", err)
		}
		fmt.Println(txName, "update account 1")
		err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.FromAccountID,
			Balance: acc1.Balance - arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("cannot update account: %w", err)
		}
		fmt.Println(txName, "get account 2")
		acc2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
		if err != nil {
			return fmt.Errorf("cannot get account: %w", err)
		}
		fmt.Println(txName, "update account 2")
		err = q.UpdateAccount(ctx, UpdateAccountParams{
			ID:      arg.ToAccountID,
			Balance: acc2.Balance + arg.Amount,
		})
		if err != nil {
			return fmt.Errorf("cannot update account: %w", err)
		}

		result.FromAccount, _ = q.GetAccount(ctx, arg.FromAccountID)
		result.ToAccount, _ = q.GetAccount(ctx, arg.ToAccountID)

		return nil
	})

	return result, err
}
