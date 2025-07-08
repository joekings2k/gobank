package db

import (
	"context"
	"testing"
	"time"

	"github.com/joekings2k/gobank/util"
	"github.com/stretchr/testify/require"
)

func createTwoRandomAccounts(t *testing.T)(Account,Account){
	account1:= createRandomAccount(t)
	account2:= createRandomAccount(t)
	return account1,account2
}
func createRandomTransfer(t *testing.T)Transfer{
	account1,account2 :=createTwoRandomAccounts(t)
	arg:= CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Amount: util.RandomMoney(),
	}
	transfer,err := testQueries.CreateTransfer(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,transfer)
	require.Equal(t , arg.FromAccountID ,transfer.FromAccountID)
	require.Equal(t , arg.ToAccountID ,transfer.ToAccountID)
	require.Equal(t , arg.Amount ,transfer.Amount)
	require.NotZero(t,transfer.ID)
	require.NotZero(t,transfer.CreatedAt)
	return transfer
}

func createMultipleTransferswithSameId(accountId1 int64, accoutId2 int64 ) Transfer{
	arg:= CreateTransferParams{
		FromAccountID: accountId1,
		ToAccountID: accoutId2,
		Amount: util.RandomMoney(),
	}
	transfer,err := testQueries.CreateTransfer(context.Background(),arg)
	if (err != nil){
		return Transfer{}
	}
	return transfer
}
func TestCreateTransfer(t *testing.T){
	
	createRandomTransfer(t)
}

func TestGetTransfer(t *testing.T){
	transfer1 := createRandomTransfer(t)
	transfer2,err := testQueries.GetTransfer(context.Background(),transfer1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,transfer2)
	require.Equal(t,transfer1.ID,transfer2.ID)
	require.Equal(t,transfer1.FromAccountID,transfer2.FromAccountID)
	require.Equal(t,transfer1.ToAccountID,transfer2.ToAccountID)
	require.Equal(t,transfer1.Amount,transfer2.Amount)
	require.WithinDuration(t,transfer1.CreatedAt.Time,transfer2.CreatedAt.Time, time.Second)
}

func TestListTransfers(t *testing.T){
	account1,account2 :=createTwoRandomAccounts(t)
	for i := 0;i <10;i ++{
		createMultipleTransferswithSameId(account1.ID,account2.ID)
	}
	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID: account2.ID,
		Limit: 5,
		Offset: 5,
	}
	transfers,err := testQueries.ListTransfers(context.Background(),arg)
	require.NoError(t,err)
	require.Len(t,transfers,5)
	for _,transfer := range transfers{
		require.NotEmpty(t,transfer)
	}
}