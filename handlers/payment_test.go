package handlers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mocks "go.simple/structure/mocks/services"
	service "go.simple/structure/services"
)

func TestCreateTransaction(t *testing.T) {
	tests := []struct {
		name                        string
		giveURI                     string
		giveReqBody                 string
		wantUserTransactionResponse *service.UserTransactionResponse
		wantErr                     error
		wantResponseBody            string
		wantHttpCode                int
	}{
		{
			name:                        "test-create-success",
			giveURI:                     "/api/users/1/transactions",
			giveReqBody:                 `{"account_id": 1, "amount": 1.0,"transaction_type": "deposit"}`,
			wantUserTransactionResponse: &service.UserTransactionResponse{1, 1, 1, 1.0, "deposit", "2020-02-10 20:10:00 +0700", ""},
			wantErr:                     nil,
			wantResponseBody:            `{"id":1,"user_id":1,"account_id":1,"amount":1,"transaction_type":"deposit","create_at":"2020-02-10 20:10:00 +0700"}`,
			wantHttpCode:                200,
		},
		{
			name:                        "test-create-error",
			giveURI:                     "/api/users/2/transactions",
			giveReqBody:                 `{"transaction_type": "deposit"}`,
			wantUserTransactionResponse: nil,
			wantErr:                     fmt.Errorf("create error"),
			wantResponseBody:            `"Internal server error"`,
			wantHttpCode:                500,
		},
		{
			name:                        "test-create-error-invalid-user_id",
			giveURI:                     "/api/users/user_id/transactions",
			giveReqBody:                 `{}`,
			wantUserTransactionResponse: nil,
			wantErr:                     nil,
			wantResponseBody:            `"Invalid user_id"`,
			wantHttpCode:                400,
		},
		{
			name:                        "test-create-error-invalid-transaction_type",
			giveURI:                     "/api/users/3/transactions",
			giveReqBody:                 `{"transaction_type": "invalid_type"}`,
			wantUserTransactionResponse: nil,
			wantErr:                     nil,
			wantResponseBody:            `"Uknown transaction type"`,
			wantHttpCode:                400,
		},
	}

	for _, v := range tests {
		// Setup service mock
		mockPs := &mocks.PaymentService{}
		mockPs.On("CreateTransaction", mock.Anything).Return(v.wantUserTransactionResponse, v.wantErr)
		ph := &PaymentHandler{paymentService: mockPs}

		w := httptest.NewRecorder()
		r := gin.Default()
		r.POST("/api/users/:user_id/transactions", ph.CreateTransaction)

		// Setup request data
		req, _ := http.NewRequest("POST", v.giveURI, bytes.NewBuffer([]byte(v.giveReqBody)))

		r.ServeHTTP(w, req)

		// Parse response data
		// var resultBody string
		b, _ := ioutil.ReadAll(w.Body)
		// _ = json.Unmarshal(b, &resultBody)

		// Check expectation
		require.Equal(t, v.wantResponseBody, string(b), v.name+": Create response didn't match")
		require.Equal(t, v.wantHttpCode, w.Code, v.name+": Create response didn't match")
	}
}

func TestGetUserTransaction(t *testing.T) {
	tests := []struct {
		name                         string
		giveURI                      string
		wantUserTransactionResponses []service.UserTransactionResponse
		wantErr                      error
		wantResponseBody             string
		wantHttpCode                 int
	}{
		{
			name:    "test-get-success",
			giveURI: "/api/users/1/transactions?account_id=1",
			wantUserTransactionResponses: []service.UserTransactionResponse{
				{1, 1, 1, 1.0, "deposit", "2020-02-10 20:10:00 +0700", ""},
			},
			wantErr:          nil,
			wantResponseBody: `[{"id":1,"user_id":1,"account_id":1,"amount":1,"transaction_type":"deposit","create_at":"2020-02-10 20:10:00 +0700"}]`,
			wantHttpCode:     200,
		},
		{
			name:                         "test-get-error-invalid_user_id",
			giveURI:                      "/api/users/user_id/transactions",
			wantUserTransactionResponses: nil,
			wantErr:                      nil,
			wantResponseBody:             `"Invalid user_id"`,
			wantHttpCode:                 400,
		},
		{
			name:                         "test-get-error-invalid_account_id",
			giveURI:                      "/api/users/1/transactions?account_id=-1",
			wantUserTransactionResponses: nil,
			wantErr:                      nil,
			wantResponseBody:             `"Invalid account_id"`,
			wantHttpCode:                 400,
		},
	}

	for _, v := range tests {
		// Setup service mock
		mockPs := &mocks.PaymentService{}
		mockPs.On("GetUserTransactions", mock.Anything, mock.Anything).Return(v.wantUserTransactionResponses, v.wantErr)
		ph := &PaymentHandler{paymentService: mockPs}

		w := httptest.NewRecorder()
		r := gin.Default()
		r.GET("/api/users/:user_id/transactions", ph.GetUserTransaction)

		// Setup request data
		req, _ := http.NewRequest("GET", v.giveURI, nil)

		r.ServeHTTP(w, req)

		// Parse response data
		// var resultBody string
		b, _ := ioutil.ReadAll(w.Body)
		// _ = json.Unmarshal(b, &resultBody)

		// Check expectation
		require.Equal(t, v.wantResponseBody, string(b), v.name+": Get response didn't match")
		require.Equal(t, v.wantHttpCode, w.Code, v.name+": Get response didn't match")
	}
}
