package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransfrTx(t *testing.T) {
	store := NewStore(testDb) // testDB es la referenia. a la cnn existente

	acount1 := createRandomAccount(t, "")
	acount2 := createRandomAccount(t, "")

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: acount1.ID,
				ToAccountID:   acount2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, acount1.ID, transfer.FromAccountID)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, acount1.ID, fromEntry.AccountID)

		// check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, acount1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, acount2.ID, toAccount.ID)

	}

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDb)
	account1 := createRandomAccount(t, "")
	account2 := createRandomAccount(t, "")

	n := 10
	errs := make(chan error)

	// La mitad: 1→2, la otra mitad: 2→1
	for i := 0; i < n; i++ {
		fromID := account1.ID
		toID := account2.ID

		if i%2 == 1 {
			fromID = account2.ID
			toID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromID,
				ToAccountID:   toID,
				Amount:        10,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}
}
