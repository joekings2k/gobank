package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockdb "github.com/joekings2k/gobank/db/mock"
	db "github.com/joekings2k/gobank/db/sqlc"
	"github.com/joekings2k/gobank/token"
	"github.com/joekings2k/gobank/util"
	"github.com/stretchr/testify/require"
)




func TestGetAcccount(t *testing.T){
	user, _ := ramdomUser(t)
	account := randomAccount(user.Username)
	testCases := [] struct {
		name string
		accountID int64
		setupAuth 		func (t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T ,recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(),gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil)
			},
			checkResponse:func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK,recorder.Code)
				requireBodyMatchAccount(t,recorder.Body,account)
			},

		},
		{
			name: "UnauthorizedUser",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, "unathorized_user", time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(),gomock.Eq(account.ID)).
				Times(1).
				Return(account, nil)
			},
			checkResponse:func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized,recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(),gomock.Any()).
				Times(0)
			},
			checkResponse:func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized,recorder.Code)
			},
		},
		{
			name: "NotFound",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
								addAuthorization(t, request, tokenMaker, authorizationTypeBearer,user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(),gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse:func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound,recorder.Code)
				
			},

		},
		{
			name: "InternalError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
								addAuthorization(t, request, tokenMaker, authorizationTypeBearer,user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(),gomock.Eq(account.ID)).
				Times(1).
				Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse:func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError,recorder.Code)
			},

		},
		{
			name: "InvalidID",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
								addAuthorization(t, request, tokenMaker, authorizationTypeBearer,user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				GetAccount(gomock.Any(),gomock.Any()).
				Times(0)
			},
			checkResponse:func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest,recorder.Code)
			},
		},
	}
	for i := range testCases{
		tc := testCases[i]
		t.Run(tc.name,func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)
			
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request,err :=http.NewRequest(http.MethodGet,url,nil)
			require.NoError(t,err)

			tc.setupAuth(t,request,server.tokenMaker)
			server.router.ServeHTTP(recorder,request)
			tc.checkResponse(t,recorder)
		})
	}
}

func TestCheckHealth (t *testing.T){
	ctrl :=gomock.NewController(t)
	defer ctrl.Finish()
	store := mockdb.NewMockStore(ctrl)

	server := newTestServer(t, store)
	recorder := httptest.NewRecorder()
	url := "/"
	request,err := http.NewRequest(http.MethodGet,url,nil)
	require.NoError(t,err)
	server.router.ServeHTTP(recorder,request)
	require.Equal(t,http.StatusOK, recorder.Code)
	
}


func TestCreateAccount(t *testing.T){
	user, _ := ramdomUser(t)
	account := randomAccount(user.Username)
	testCases := [] struct {
		name string
		body gin.H
		setupAuth 		func (t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs 		func(store *mockdb.MockStore)
		checkResponse func (t *testing.T ,recorder *httptest.ResponseRecorder)

	}{
		{
			name: "OK",
			body: gin.H{
				"owner": account.Owner,
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				CreateAccount(gomock.Any(),gomock.Eq(db.CreateAccountParams{
					Owner: account.Owner,
					Currency: account.Currency,
					Balance: 0,
				})).Times(1).Return(account,nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK,recorder.Code)
				requireBodyMatchAccount(t,recorder.Body,account)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"owner": account.Owner,
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				CreateAccount(gomock.Any(),gomock.Eq(db.CreateAccountParams{
					Owner: account.Owner,
					Currency: account.Currency,
					Balance: 0,
				})).Times(1).Return(account,sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError,recorder.Code)
				
			},
		},
		{
			name: "InvalidParameters",
			body: gin.H{
				"owner": "",
				"currency": "",
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				CreateAccount(gomock.Any(),gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest,recorder.Code)
			},
		},


	}

	for i := range testCases{
		tc := testCases[i]
		t.Run(tc.name,func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store:=mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server:= newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t,err)
			url := "/accounts"
			request,err := http.NewRequest(http.MethodPost,url,bytes.NewReader(data))
			require.NoError(t,err)
			tc.setupAuth(t,request,server.tokenMaker)
			request.Header.Set("Content-Type", "application/json")
			server.router.ServeHTTP(recorder,request)
			tc.checkResponse(t,recorder)
		})
	}
}

func TestListAccounts(t *testing.T){
	user, _ := ramdomUser(t)
	accounts :=make([]db.Account,10)
	for i:=0;i <10 ;i++{
		accounts[i] = randomAccount(user.Username)
	}
	testCases := []struct{
		name string
		parms ListAccountsRequest
		setupAuth 		func (t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs 		func(store *mockdb.MockStore)
		checkResponse func (t *testing.T ,recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			parms: ListAccountsRequest{
				PageID: 1,
				PageSize: 5,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker,authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(),gomock.Eq(db.ListAccountsParams{
					Owner: user.Username,
					Limit:5,
					Offset: (0)*5,
				})).Times(1).Return(accounts,nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusOK,recorder.Code)
				requireBodyMatchAccounts(t,recorder.Body,accounts)
			},
		},
		{
			name: "InvalidParameters",
			parms: ListAccountsRequest{
				PageID: 0,
				PageSize: 5,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(),gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusBadRequest,recorder.Code)
			},
		},
		{
			name: "InternalError",
			parms: ListAccountsRequest{
				PageID: 1,
				PageSize: 5,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(),gomock.Eq(db.ListAccountsParams{
					Owner: user.Username ,
					Limit:5,
					Offset: (0)*5,
				})).Times(1).Return([]db.Account{},sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusInternalServerError,recorder.Code)
			},
		},
	}
	for i:= range testCases{
		tc :=testCases[i]
		t.Run(tc.name,func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d",tc.parms.PageID,tc.parms.PageSize)
			request,err :=http.NewRequest(http.MethodGet,url,nil)
			require.NoError(t,err)
			tc.setupAuth(t,request,server.tokenMaker)
			server.router.ServeHTTP(recorder,request)
			tc.checkResponse(t,recorder)
		})
	}
}

func TestUpdateAccount(t *testing.T){
	user, _ := ramdomUser(t)
	account := randomAccount(user.Username)
	updatedAccount := account
	updatedAccount.Balance = 100 
	testCases := []struct{
		name string
		accountID int64
		body updateAccountBody
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T,recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			accountID: account.ID,
			body: updateAccountBody{
				Balance: 100,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAccount(gomock.Any(),gomock.Eq(db.UpdateAccountParams{
					ID: account.ID,
					Balance: 100,
				})).Times(1).Return(updatedAccount,nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusOK,recorder.Code)
				requireBodyMatchAccount(t,recorder.Body,updatedAccount)
			},
			
		},
		{
			name: "InvalidID",
			accountID: 0,
			body: updateAccountBody{
				Balance: 100,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAccount(gomock.Any(),gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusBadRequest,recorder.Code)
				
			},
		},
		{
			name: "InvalidBalance",
			accountID: account.ID,
			body: updateAccountBody{},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAccount(gomock.Any(),gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusBadRequest,recorder.Code)
				
			},
		},
		{
			name: "InternalError",
			accountID: account.ID,
			body: updateAccountBody{
				Balance: 100,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().UpdateAccount(gomock.Any(),gomock.Eq(db.UpdateAccountParams{
					ID: account.ID,
					Balance: 100,
				})).Times(1).Return(db.Account{},sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t,http.StatusInternalServerError,recorder.Code)
			},
		},
	}
	for i := range testCases{
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			
			
			data,err := json.Marshal(tc.body)
			require.NoError(t,err)
			url := fmt.Sprintf("/accounts/%d",tc.accountID)
			request,err := http.NewRequest(http.MethodPatch,url,bytes.NewReader(data))
			require.NoError(t,err)
			request.Header.Set("Content-Type", "application/json")
			server.router.ServeHTTP(recorder,request)
			tc.checkResponse(t,recorder)
		})
	}
}


func TestDeleteAccount(t *testing.T) {
	user, _ := ramdomUser(t)
	account := randomAccount(user.Username)
	testCases := []struct{
		name string
		accountID int64
		buildStubs func (store *mockdb.MockStore)
		checkResponse func (t *testing.T , recorder *httptest.ResponseRecorder)
	}{
		{
			name : "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				DeleteAccount(gomock.Any(),gomock.Eq(account.ID)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK,recorder.Code)
			},
		},
		{
			name : "InvalidID",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				DeleteAccount(gomock.Any(),gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest,recorder.Code)
			},
		},
		{
			name : "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
				DeleteAccount(gomock.Any(),gomock.Eq(account.ID)).Times(1).Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError,recorder.Code)
			},
		},
	}
	for i :=range testCases{
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			// create controller
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			//create store with controller
			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			// create server
			server := newTestServer(t,store)
			// create recorder
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request,err :=http.NewRequest(http.MethodDelete,url,nil)
			require.NoError(t,err)
			server.router.ServeHTTP(recorder,request)
			tc.checkResponse(t,recorder)

		})
	}

}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:util.RandomInt(1,1000),
		Owner: owner,
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data,err := io.ReadAll(body)
	require.NoError(t,err )

	var gotAccount  db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t,err)
	require.Equal(t,account,gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}