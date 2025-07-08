package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/joekings2k/gobank/util"
	"github.com/stretchr/testify/require"
)


func createRandomEntry (t *testing.T)Entry{
	acc := createRandomAccount(t)
	arg := CreateEntryParams {
		AccountID:acc.ID ,
		Amount: util.RandomMoney(),
	}
	entry,err := testQueries.CreateEntry(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,entry)
	require.Equal(t,arg.AccountID,entry.AccountID)
	require.Equal(t,arg.Amount,entry.Amount)
	require.NotZero(t,entry.ID)
	require.NotZero(t,entry.CreatedAt)

return entry
}

func createRandomEntryWithSameAccountId(accountId int64) Entry{
	arg := CreateEntryParams {
		AccountID: accountId ,
		Amount: util.RandomMoney(),
	}
	entry,err := testQueries.CreateEntry(context.Background(),arg)
	if (err != nil){
		return Entry{}
	}
	return entry
}

func TestCreateEntry (t *testing.T){
	createRandomEntry(t)
}

func TestDeleteEntry (t *testing.T){
	entry1 := createRandomEntry(t)
	err :=testQueries.DeleteEntry(context.Background(),entry1.ID)
	require.NoError(t,err)
	entry2,err := testQueries.GetEntry(context.Background(),entry1.ID)
	require.Error(t,err)
	require.EqualError(t,err,sql.ErrNoRows.Error())
	require.Empty(t,entry2)
}

func TestGetEntry (t *testing.T){
	entry1 := createRandomEntry(t)
	entry2,err := testQueries.GetEntry(context.Background(),entry1.ID)
	require.NoError(t,err)
	require.NotEmpty(t,entry2)
	require.Equal(t,entry1.ID,entry2.ID)
	require.Equal(t,entry1.AccountID,entry2.AccountID)
	require.Equal(t,entry1.Amount,entry2.Amount)

	require.WithinDuration(t,entry1.CreatedAt,entry2.CreatedAt,time.Second)
}

func TestListEntries (t *testing.T){
	acc := createRandomAccount(t)

	for i := 0;i <10;i ++{
		createRandomEntryWithSameAccountId(acc.ID)
	}
	arg := ListEntriesParams{
		AccountID: acc.ID,
		Limit: 5,
		Offset: 5,
	}
	entries,err := testQueries.ListEntries(context.Background(),arg)
	require.NoError(t,err)
	require.Len(t,entries,5)
	for _,entry := range entries{
		require.NotEmpty(t,entry)
	}
}