package repositories

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupRatingProductRepoTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *tsImpl) {
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialect := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.NoError(t, err)

	ts := tsImpl{
		db: db,
	}
	return sqlDB, mock, &ts
}

func TestCreateTransaction(t *testing.T) {
	tests := []struct {
		giveTS          Transaction
		wantRowAffected int64
		wantErr         error
	}{
		{Transaction{0, 1, 1, 1.0, "type1", time.Now(), time.Now()}, 1, nil},
		{Transaction{0, 2, 2, 2.0, "type2", time.Time{}, time.Now()}, 0, fmt.Errorf("create error")},
	}

	sqlDB, mock, ts := setupRatingProductRepoTest(t)
	defer sqlDB.Close()

	for _, v := range tests {
		// Setup mock scenario
		mock.ExpectBegin()
		if v.wantRowAffected == 0 {
			mock.ExpectExec("INSERT INTO `transactions`").
				WithArgs(v.giveTS.UserID, v.giveTS.AccountID, v.giveTS.Amount, v.giveTS.Type, v.giveTS.CreateAt, v.giveTS.UpdateAt).
				WillReturnError(v.wantErr)
			mock.ExpectRollback()
		} else {
			mock.ExpectExec("INSERT INTO `transactions`").
				WithArgs(v.giveTS.UserID, v.giveTS.AccountID, v.giveTS.Amount, v.giveTS.Type, v.giveTS.CreateAt, v.giveTS.UpdateAt).
				WillReturnResult(sqlmock.NewResult(1, v.wantRowAffected))
			mock.ExpectCommit()
		}

		result, err := ts.CreateTransaction(v.giveTS)

		// Check Expectations
		require.True(t, mock.ExpectationsWereMet() == nil, "Mock didn't meet expectation")
		require.Equal(t, v.wantRowAffected, result, "Rows affected didn't match")
		require.Equal(t, v.wantErr, err, "Return error didnt match")
	}
}

func TestGetUserTransactionsWithAccount(t *testing.T) {
	tests := struct {
		giveUserID       int64
		giveAccountID    int64
		wantTransactions []Transaction
		wantErr          error
	}{
		1, 2, []Transaction{{1, 1, 1, 1.0, "type1", time.Now(), time.Now()}}, nil,
	}
	sqlDB, mock, ts := setupRatingProductRepoTest(t)
	defer sqlDB.Close()

	var query string
	cols := []string{"id", "user_id", "account_id", "amount", "type", "create_at", "update_at"}

	// Setup mock scenario
	rows := sqlmock.NewRows(cols).AddRow(tests.wantTransactions[0].ID, tests.wantTransactions[0].UserID, tests.wantTransactions[0].AccountID,
		tests.wantTransactions[0].Amount, tests.wantTransactions[0].Type, tests.wantTransactions[0].CreateAt, tests.wantTransactions[0].UpdateAt)
	query = "SELECT * FROM `transactions` WHERE user_id = ? AND account_id = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(tests.giveUserID, tests.giveAccountID).
		WillReturnRows(rows)

	result, err := ts.GetUserTransactions(tests.giveUserID, tests.giveAccountID)

	// Check Expectations
	require.True(t, mock.ExpectationsWereMet() == nil, "Mock didn't meet expectation")
	require.Equal(t, tests.wantTransactions, result, "Transaction response didn't match")
	require.Equal(t, tests.wantErr, err, "Return error didnt match")
}

func TestGetUserTransactionsNoAccount(t *testing.T) {
	tests := struct {
		giveUserID       int64
		giveAccountID    int64
		wantTransactions []Transaction
		wantErr          error
	}{
		2, 0, []Transaction{{1, 1, 1, 1.0, "type1", time.Now(), time.Now()}, {2, 2, 2, 2.0, "type2", time.Now(), time.Now()}}, nil,
	}

	sqlDB, mock, ts := setupRatingProductRepoTest(t)
	defer sqlDB.Close()

	var query string
	cols := []string{"id", "user_id", "account_id", "amount", "type", "create_at", "update_at"}

	// Setup mock scenario
	rows := sqlmock.NewRows(cols).AddRow(tests.wantTransactions[0].ID, tests.wantTransactions[0].UserID, tests.wantTransactions[0].AccountID,
		tests.wantTransactions[0].Amount, tests.wantTransactions[0].Type, tests.wantTransactions[0].CreateAt, tests.wantTransactions[0].UpdateAt).
		AddRow(tests.wantTransactions[1].ID, tests.wantTransactions[1].UserID, tests.wantTransactions[1].AccountID,
			tests.wantTransactions[1].Amount, tests.wantTransactions[1].Type, tests.wantTransactions[1].CreateAt, tests.wantTransactions[1].UpdateAt)
	query = "SELECT * FROM `transactions` WHERE user_id = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(tests.giveUserID).
		WillReturnRows(rows)

	result, err := ts.GetUserTransactions(tests.giveUserID, tests.giveAccountID)

	// Check Expectations
	require.True(t, mock.ExpectationsWereMet() == nil, "Mock didn't meet expectation")
	require.Equal(t, tests.wantTransactions, result, "Transaction response didn't match")
	require.Equal(t, tests.wantErr, err, "Return error didnt match")
}

func TestGetUserTransactionsOnFailure(t *testing.T) {
	tests := struct {
		giveUserID       int64
		giveAccountID    int64
		wantTransactions []Transaction
		wantErr          error
	}{
		1, 2, nil, fmt.Errorf("error get"),
	}

	sqlDB, mock, ts := setupRatingProductRepoTest(t)
	defer sqlDB.Close()

	// Setup mock scenario
	query := "SELECT * FROM `transactions` WHERE user_id = ? AND account_id = ?"
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(tests.giveUserID, tests.giveAccountID).
		WillReturnError(tests.wantErr)

	result, err := ts.GetUserTransactions(tests.giveUserID, tests.giveAccountID)

	// Check Expectations
	require.True(t, mock.ExpectationsWereMet() == nil, "Mock didn't meet expectation")
	require.Equal(t, tests.wantTransactions, result, "Transaction response didn't match")
	require.Equal(t, tests.wantErr, err, "Return error didnt match")
}
