package service

import (
	"fmt"
	"time"

	"go.simple/structure/repositories"
	"go.simple/structure/utils"
	"gorm.io/gorm"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type PaymentService interface {
	CreateTransaction(req CreateTransactionRequest) (*UserTransactionResponse, error)
	GetUserTransactions(userID, accountID int64) ([]UserTransactionResponse, error)
	UpdateUserTransaction(req UpdateTransactionRequest) (*UserTransactionResponse, error)
	DeleteUserTransaction(id, userID int64) error
}

type psImpl struct {
	transactionRepo repositories.TransactionRepo
}

func NewPaymentService(db *gorm.DB) PaymentService {
	return &psImpl{
		transactionRepo: repositories.NewTransactionRepo(db),
	}
}

type CreateTransactionRequest struct {
	UserID          int64   `json:"user_id"`
	AccountID       int64   `json:"account_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
}

type UserTransactionResponse struct {
	ID              int64   `json:"id"`
	UserID          int64   `json:"user_id"`
	AccountID       int64   `json:"account_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	CreateAt        string  `json:"create_at"`
	UpdateAt        string  `json:"update_at,omitempty"`
}

func (ps *psImpl) CreateTransaction(req CreateTransactionRequest) (*UserTransactionResponse, error) {
	tsData := repositories.Transaction{
		UserID:    req.UserID,
		AccountID: req.AccountID,
		Amount:    req.Amount,
		Type:      req.TransactionType,
		CreateAt:  time.Now().UTC(),
		UpdateAt:  time.Now().UTC(),
	}

	id, err := ps.transactionRepo.CreateTransaction(tsData)
	if err != nil {
		return nil, err
	}

	result := UserTransactionResponse{
		ID:              id,
		UserID:          req.UserID,
		AccountID:       req.AccountID,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		CreateAt:        utils.FormatVNTime(tsData.CreateAt, tsCreateTimeLayout),
	}
	return &result, nil
}

func (ps *psImpl) GetUserTransactions(userID, accountID int64) ([]UserTransactionResponse, error) {
	userTs, err := ps.transactionRepo.GetUserTransactions(userID, accountID)
	if err != nil {
		return nil, err
	}

	var result []UserTransactionResponse
	for _, v := range userTs {
		result = append(result, UserTransactionResponse{
			ID:              v.ID,
			UserID:          v.UserID,
			AccountID:       v.AccountID,
			Amount:          v.Amount,
			TransactionType: v.Type,
			CreateAt:        utils.FormatVNTime(v.CreateAt, tsCreateTimeLayout),
			UpdateAt:        utils.FormatVNTime(v.UpdateAt, tsCreateTimeLayout),
		})
	}

	return result, nil
}

type UpdateTransactionRequest struct {
	ID              int64   `json:"id"`
	UserID          int64   `json:"user_id"`
	AccountID       int64   `json:"account_id"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
}

var ZeroAffectedErr = fmt.Errorf("Zero affected rows")

func (ps *psImpl) UpdateUserTransaction(req UpdateTransactionRequest) (*UserTransactionResponse, error) {
	tsData := repositories.Transaction{
		ID:        req.ID,
		UserID:    req.UserID,
		AccountID: req.AccountID,
		Amount:    req.Amount,
		Type:      req.TransactionType,
		UpdateAt:  time.Now().UTC(),
	}

	rowsAffected, err := ps.transactionRepo.UpdateUserTransactions(tsData)
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ZeroAffectedErr
	}

	result := UserTransactionResponse{
		ID:              req.ID,
		UserID:          req.UserID,
		AccountID:       req.AccountID,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		CreateAt:        utils.FormatVNTime(tsData.UpdateAt, tsCreateTimeLayout),
	}
	return &result, nil
}

func (ps *psImpl) DeleteUserTransaction(id, userID int64) error {
	rowsAffected, err := ps.transactionRepo.DeleteUserTransactions(id, userID)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ZeroAffectedErr
	}
	return nil
}
