package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	mocks "go.simple/structure/mocks/repositories"
	"go.simple/structure/repositories"
	"go.simple/structure/utils"
)

func TestCreateTransaction(t *testing.T) {
	tests := []struct {
		name                     string
		giveTransactionRequest   CreateTransactionRequest
		wantTransactionsResponse *UserTransactionResponse
		wantRowAffected          int64
		wantErr                  error
	}{
		{
			name:                     "test-successs",
			giveTransactionRequest:   CreateTransactionRequest{1, 1, 1, "type1"},
			wantTransactionsResponse: &UserTransactionResponse{1, 1, 1, 1, "type1", utils.FormatVNTime(time.Now().UTC(), tsCreateTimeLayout), ""},
			wantRowAffected:          1,
			wantErr:                  nil,
		},
		{
			name:                     "test-error",
			giveTransactionRequest:   CreateTransactionRequest{1, 1, 1, "type2"},
			wantTransactionsResponse: nil,
			wantRowAffected:          0,
			wantErr:                  fmt.Errorf("create error"),
		},
	}

	for _, v := range tests {
		mockTs := &mocks.TransactionRepo{}
		ps := psImpl{transactionRepo: mockTs}
		mockTs.On("CreateTransaction", mock.Anything).Return(v.wantRowAffected, v.wantErr)

		result, err := ps.CreateTransaction(v.giveTransactionRequest)

		// Check Expectations
		require.Equal(t, v.wantTransactionsResponse, result, v.name+": Transactions response didn't match")
		require.Equal(t, v.wantErr, err, v.name+": Return error didnt match")
	}
}

func TestGetUserTransactions(t *testing.T) {
	tests := []struct {
		name                          string
		giveUserID                    int64
		giveAccountID                 int64
		wantTransaction               []repositories.Transaction
		wantUserTransactionsResponses []UserTransactionResponse
		wantErr                       error
	}{
		{
			name:          "test-successs",
			giveUserID:    1,
			giveAccountID: 1,
			wantTransaction: []repositories.Transaction{
				{1, 1, 1, 1.0, "type1", time.Now().UTC(), time.Time{}},
			},
			wantUserTransactionsResponses: []UserTransactionResponse{
				{1, 1, 1, 1.0, "type1", utils.FormatVNTime(time.Now().UTC(), tsCreateTimeLayout), ""},
			},
			wantErr: nil,
		},
		{
			name:          "test-successs-have-update_at",
			giveUserID:    1,
			giveAccountID: 0,
			wantTransaction: []repositories.Transaction{
				{1, 1, 0, 1.0, "type1", time.Now().UTC(), time.Now().UTC()},
			},
			wantUserTransactionsResponses: []UserTransactionResponse{
				{1, 1, 0, 1.0, "type1", utils.FormatVNTime(time.Now().UTC(), tsCreateTimeLayout), utils.FormatVNTime(time.Now().UTC(), tsCreateTimeLayout)},
			},
			wantErr: nil,
		},
		{
			name:                          "test-error",
			giveUserID:                    1,
			giveAccountID:                 1,
			wantTransaction:               nil,
			wantUserTransactionsResponses: nil,
			wantErr:                       fmt.Errorf("get error"),
		},
	}

	for _, v := range tests {
		mockTs := &mocks.TransactionRepo{}
		ps := psImpl{transactionRepo: mockTs}
		mockTs.On("GetUserTransactions", mock.Anything, mock.Anything).Return(v.wantTransaction, v.wantErr)

		result, err := ps.GetUserTransactions(v.giveUserID, v.giveAccountID)

		// Check Expectations
		require.Equal(t, v.wantUserTransactionsResponses, result, v.name+": Transactions response didn't match")
		require.Equal(t, v.wantErr, err, v.name+": Return error didnt match")
	}
}
