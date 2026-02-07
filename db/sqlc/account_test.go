package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valvarez/simplebank/utils"
)

func createRandomAccount(t *testing.T, userName string) Account {

	if userName == "" {
		userName = utils.RandomOwner()
	}

	arg := CreateAccountsParams{
		Owner:    userName,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQuerys.CreateAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	account1 := createRandomAccount(t, "")

	err := testQuerys.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
}

func TestGetAccount(t *testing.T) {

	account1 := createRandomAccount(t, "")

	account2, err := testQuerys.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, 0)

	err = testQuerys.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t, "")

	err := testQuerys.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testQuerys.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, "sql: no rows in result set")
	require.Equal(t, sql.ErrNoRows, err)
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	createRandomAccount(t, "")
	accounts, err := testQuerys.ListAccounts(context.Background())
	require.NoError(t, err)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

	//delete accounts
	for _, account := range accounts {
		err := testQuerys.DeleteAccount(context.Background(), account.ID)
		require.NoError(t, err)
	}
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t, "")

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2, err := testQuerys.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, arg.Balance, account2.Balance)
	require.NotEqual(t, account1.Balance, account2.Balance)

	err = testQuerys.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)
}
