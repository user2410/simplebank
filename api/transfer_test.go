package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "github.com/user2410/simplebank/db/mock"
	db "github.com/user2410/simplebank/db/sqlc"
	"github.com/user2410/simplebank/util"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransferAPI(t *testing.T) {
	amount := int64(10)

	user1, _ := randomUser(t)
	user2, _ := randomUser(t)
	user3, _ := randomUser(t)

	acc1 := randomAccount(user1.Username)
	acc2 := randomAccount(user2.Username)
	acc3 := randomAccount(user3.Username)

	acc1.Currency = util.Currencies[0]
	acc2.Currency = util.Currencies[0]
	acc3.Currency = util.Currencies[1]

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc2.ID,
				"amount":          amount,
				"currency":        util.Currencies[0],
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).Times(1).Return(acc1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc2.ID)).Times(1).Return(acc2, nil)

				arg := db.TransferTxParams{
					FromAccountID: acc1.ID,
					ToAccountID:   acc2.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(arg)).Times(1)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "InvalidCurrency",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc2.ID,
				"amount":          amount,
				"currency":        "XYZ",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "NegativeAmount",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc2.ID,
				"amount":          -amount,
				"currency":        util.Currencies[0],
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "GetAccountError",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc2.ID,
				"amount":          amount,
				"currency":        util.Currencies[0],
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "TransferTxError",
			body: gin.H{
				"from_account_id": acc1.ID,
				"to_account_id":   acc2.ID,
				"amount":          amount,
				"currency":        util.Currencies[0],
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc1.ID)).Times(1).Return(acc1, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(acc2.ID)).Times(1).Return(acc2, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(1).Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			store := mockdb.NewMockStore(ctl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/transfers", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
