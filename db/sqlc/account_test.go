package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"simplebank/util"
	"testing"
	"time"

	_ "github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	args := CreateAccountParams{
		Owner:       util.RandomOwner(),
		Balance:     util.RandomMoney(),
		Currency:    util.RandomCurrency(),
		CountryCode: util.RandomCountryCode(),
	}
	acc, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.Equal(t, args.Owner, acc.Owner)
	require.Equal(t, args.Balance, acc.Balance)
	require.Equal(t, args.Currency, acc.Currency)
	require.NotZero(t, acc.ID)
	require.NotZero(t, acc.CreatedAt)
	return acc
}

func TestQueries_CreateAccount(t *testing.T) {
	acc := createRandomAccount(t)

	// cleanup
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc.ID))
}

func TestQueries_GetAccount(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2, err := testQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, acc2)
	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, acc1.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.Equal(t, acc1.CountryCode, acc2.CountryCode)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)

	// cleanup
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc1.ID))
}

func TestQueries_UpdateAccount(t *testing.T) {
	acc1 := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:      acc1.ID,
		Balance: util.RandomMoney(),
	}

	acc2, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, acc2)
	require.Equal(t, acc1.ID, acc2.ID)
	require.Equal(t, acc1.Owner, acc2.Owner)
	require.Equal(t, args.Balance, acc2.Balance)
	require.Equal(t, acc1.Currency, acc2.Currency)
	require.Equal(t, acc1.CountryCode, acc2.CountryCode)
	require.WithinDuration(t, acc1.CreatedAt, acc2.CreatedAt, time.Second)

	// cleanup
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc1.ID))
}

func TestQueries_DeleteAccount(t *testing.T) {
	acc := createRandomAccount(t)
	require.NoError(t, testQueries.DeleteAccount(context.Background(), acc.ID))

	acc2, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}
