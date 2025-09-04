package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/joekings2k/gobank/db/mock"
	db "github.com/joekings2k/gobank/db/sqlc"
	"github.com/joekings2k/gobank/util"
	"github.com/stretchr/testify/require"
)

func  TestCreateTransfer(t *testing.T) {
	account1 := randomAccount()
	account2 := randomAccount()
	account3 := randomAccount()

	account1.Currency = util.USD
	account2.Currency = util.USD
	account3.Currency = util.EUR
	testCases := [] struct {
		name string
		body gin.H
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name :"OK",
			body: gin.H{
				"from_account_id":account1.ID,
				"to_account_id":account2.ID,
				"amount":10,
				"currency":util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account1.ID)).Times(1).Return(account1,nil)
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account2.ID)).Times(1).Return(account2,nil)
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID: account2.ID,
					Amount: 10,
				}
				store.EXPECT().TransferTx(gomock.Any(),gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusOK,recorder.Code)
			},
		},
		{
			name :"InvalidAccountID",
			body: gin.H{
				"from_account_id":0,
				"to_account_id":0,
				"amount":10,
				"currency":util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(),gomock.Any()).Times(0)
				store.EXPECT().GetAccount(gomock.Any(),gomock.Any()).Times(0)
				arg := db.TransferTxParams{
					FromAccountID: 0,
					ToAccountID: 0,
					Amount: 10,
				}
				store.EXPECT().TransferTx(gomock.Any(),gomock.Eq(arg)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusBadRequest,recorder.Code)
			},
		},
		{
			name :"NotFound",
			body: gin.H{
				"from_account_id":account1.ID,
				"to_account_id":account2.ID,
				"amount":10,
				"currency":util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account1.ID)).Times(1).Return(db.Account{},sql.ErrNoRows)
				
				
				store.EXPECT().TransferTx(gomock.Any(),gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusNotFound,recorder.Code)
			},
		},
		{
			name :"InternalError",
			body: gin.H{
				"from_account_id":account1.ID,
				"to_account_id":account2.ID,
				"amount":10,
				"currency":util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account1.ID)).Times(1).Return(db.Account{},sql.ErrConnDone)
				
				
				store.EXPECT().TransferTx(gomock.Any(),gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusInternalServerError,recorder.Code)
			},
		},
		{
			name :"CurrencyMismatch",
			body: gin.H{
				"from_account_id":account1.ID,
				"to_account_id":account3.ID,
				"amount":10,
				"currency":util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account1.ID)).Times(1).Return(account1,nil)
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account3.ID)).Times(1).Return(account3,nil)
				
				store.EXPECT().TransferTx(gomock.Any(),gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusBadRequest,recorder.Code)
			},
		},
		{
			name :"TransferServerError",
			body: gin.H{
				"from_account_id":account1.ID,
				"to_account_id":account2.ID,
				"amount":10,
				"currency":util.USD,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account1.ID)).Times(1).Return(account1,nil)
				store.EXPECT().GetAccount(gomock.Any(),gomock.Eq(account2.ID)).Times(1).Return(account2,nil)
				
				arg := db.TransferTxParams{
					FromAccountID: account1.ID,
					ToAccountID: account2.ID,
					Amount: 10,
				}
				store.EXPECT().TransferTx(gomock.Any(),gomock.Eq(arg)).Times(1).Return(db.TransferTxResult{},sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusInternalServerError,recorder.Code)
			},
		},
	}
	for i  := range testCases {
		tc := testCases[i]
		t.Run(tc.name,func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store :=mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server:= NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/transfers"
			data,err := json.Marshal(tc.body)
			require.NoError(t,err)
			request ,err := http.NewRequest(http.MethodPost,url,bytes.NewReader(data))
			require.NoError(t,err)
			request.Header.Set("Content-Type", "application/json")
			server.router.ServeHTTP(recorder,request)
			tc.checkResponse(t,recorder)
		})
	}
	
}