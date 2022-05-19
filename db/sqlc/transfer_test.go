package db

import (
	"context"
	"testing"

	"github.com/cyberdr0id/bank-backend/util"
	"github.com/stretchr/testify/require"
)

func createRandomTransfer(t *testing.T, acc1, acc2 Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc2.ID,
		Amount:        util.RandMoney(),
	}

	result, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, result)

	id, _ := result.LastInsertId()
	require.NotZero(t, id)

	return Transfer{
		ID:            id,
		FromAccountID: arg.FromAccountID,
		ToAccountID:   arg.ToAccountID,
		Amount:        arg.Amount,
	}
}

func TestCreateTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	createRandomTransfer(t, acc1, acc2)
}

func TestGetTransfer(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	transfer := createRandomTransfer(t, acc1, acc2)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer.ID, transfer2.ID)
	require.Equal(t, transfer.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer.Amount, transfer2.Amount)
}

func TestListTransfers(t *testing.T) {
	acc1 := createRandomAccount(t)
	acc2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTransfer(t, acc1, acc2)
		createRandomTransfer(t, acc2, acc1)
	}

	arg := ListTransfersParams{
		FromAccountID: acc1.ID,
		ToAccountID:   acc1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)

		require.True(t, transfer.FromAccountID == acc1.ID || transfer.ToAccountID == acc1.ID)
	}
}
