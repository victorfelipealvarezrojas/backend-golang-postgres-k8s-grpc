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
		db:      db,      // ← Guarda la referencia a una conexión YA EXISTENTE
		Queries: New(db), // ← Crea Queries usando esa conexión existente
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // genera la transaccion
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
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

func (store *Store) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		// transferenia inicial
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})

		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})

		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})

		if err != nil {
			return err
		}

		if args.FromAccountID < args.ToAccountID {
			// Update accounts' balances
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     args.FromAccountID,
				Amount: -args.Amount,
			})

			if err != nil {
				return err
			}

			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     args.ToAccountID,
				Amount: args.Amount,
			})

			if err != nil {
				return err
			}
		} else {

			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     args.ToAccountID,
				Amount: args.Amount,
			})

			if err != nil {
				return err
			}

			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     args.FromAccountID,
				Amount: -args.Amount,
			})

			if err != nil {
				return err
			}

		}

		return nil
	})

	return result, err

}
