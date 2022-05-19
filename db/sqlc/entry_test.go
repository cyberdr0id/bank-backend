package db

import (
	"context"
	"testing"

	"github.com/cyberdr0id/bank-backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T, a Account) Entry {
	arg := CreateEntryParams{
		AccountID: a.ID,
		Amount:    util.RandMoney(),
	}

	result, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	id, _ := result.LastInsertId()
	require.NotZero(t, id)

	return Entry{
		ID:        id,
		AccountID: arg.AccountID,
		Amount:    arg.Amount,
	}
}

func TestCreateEntry(t *testing.T) {
	acc := createRandomAccount(t)
	createRandomEntry(t, acc)
}

func TestGetEntry(t *testing.T) {
	acc := createRandomAccount(t)
	entry1 := createRandomEntry(t, acc)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
}

func TestListEntries(t *testing.T) {
	acc := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntry(t, acc)
	}

	arg := ListEntriesParams{
		AccountID: acc.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
