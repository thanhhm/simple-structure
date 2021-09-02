package repositories

import (
	"time"

	"gorm.io/gorm"
)

type TransactionRepo interface {
	CreateTransaction(tsInfo Transaction) (int64, error)
	GetUserTransactions(userID, accountID int64) ([]Transaction, error)
	UpdateUserTransactions(tsInfo Transaction) (int64, error)
	DeleteUserTransactions(id, userID int64) (int64, error)
}

type tsImpl struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) TransactionRepo {
	return &tsImpl{
		db: db,
	}
}

type Transaction struct {
	ID        int64     `gorm:"id"`
	UserID    int64     `gorm:"user_id"`
	AccountID int64     `gorm:"account_id"`
	Amount    float64   `gorm:"amount_id"`
	Type      string    `gorm:"type_id"`
	CreateAt  time.Time `gorm:"create_at"`
	UpdateAt  time.Time `gorm:"update_at"`
}

func (ts *tsImpl) CreateTransaction(tsInfo Transaction) (int64, error) {
	if err := ts.db.Create(&tsInfo).Error; err != nil {
		return 0, err
	}

	return tsInfo.ID, nil
}

func (ts *tsImpl) GetUserTransactions(userID, accountID int64) ([]Transaction, error) {
	var result []Transaction
	var err error

	if accountID > 0 {
		err = ts.db.Where("user_id = ? AND account_id = ?", userID, accountID).Find(&result).Error
	} else {
		err = ts.db.Where("user_id = ?", userID).Find(&result).Error
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ts *tsImpl) UpdateUserTransactions(tsInfo Transaction) (int64, error) {
	result := ts.db.Model(&Transaction{}).Where("id=? AND user_id=?", tsInfo.ID, tsInfo.UserID).Updates(tsInfo)
	return result.RowsAffected, result.Error
}

func (ts *tsImpl) DeleteUserTransactions(id, userID int64) (int64, error) {
	result := ts.db.Where("id=? AND user_id=?", id, userID).Delete(&Transaction{})
	return result.RowsAffected, result.Error
}
