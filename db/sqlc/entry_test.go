package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valvarez/simplebank/utils"
)

func TestCreateEntry(t *testing.T) {

	entry := CreateEntryParams{
		AccountID: createRandomAccount(t, "").ID,
		Amount:    utils.RandomMoney(),
	}

	createdEntry, err := testQuerys.CreateEntry(context.Background(), entry)

	require.NoError(t, err)
	require.NotEmpty(t, createdEntry)

	require.Equal(t, entry.AccountID, createdEntry.AccountID)
	require.Equal(t, entry.Amount, createdEntry.Amount)

	err = testQuerys.DeleteEntry(context.Background(), createdEntry.ID)
	require.NoError(t, err)

	err = testQuerys.DeleteAccount(context.Background(), entry.AccountID)
	require.NoError(t, err)

}
