package db

import (
	"context"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	conn := ConnectDB(t)
	store := NewStore(conn)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// Run n concurrent transfer transactions
	n := 2
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()

	}

	// Check results
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err, fmt.Sprint(err))

		result := <-results
		assert.NotEmpty(t, result)

		// Check transfer
		transfer := result.Transfer
		assert.NotEmpty(t, transfer)
		assert.Equal(t, account1.ID, transfer.FromAccountID)
		assert.Equal(t, account2.ID, transfer.ToAccountID)
		assert.Equal(t, amount, transfer.Amount)
		assert.NotZero(t, transfer.ID)
		assert.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTranser(context.Background(), transfer.ID)
		assert.NoError(t, err)

		// Check entries
		fromEntry := result.FromEntry
		assert.NotEmpty(t, fromEntry)
		assert.Equal(t, account1.ID, fromEntry.AccountID)
		assert.Equal(t, -amount, fromEntry.Amount)
		assert.NotZero(t, fromEntry.ID)
		assert.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		assert.NoError(t, err)

		toEntry := result.ToEntry
		assert.NotEmpty(t, toEntry)
		assert.Equal(t, account2.ID, toEntry.AccountID)
		assert.Equal(t, amount, toEntry.Amount)
		assert.NotZero(t, toEntry.ID)
		assert.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		assert.NoError(t, err)

		// check account
		fromAccount := result.FromAccount
		assert.NotEmpty(t, fromAccount)
		assert.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		assert.NotEmpty(t, toAccount)
		assert.Equal(t, toAccount.ID, account2.ID)

		// // Check account balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		assert.Equal(t, diff1, diff2)
		assert.True(t, diff1 > 0)
		assert.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		assert.True(t, k >= 1 && k <= n)
		assert.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the final update balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)

	assert.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	assert.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	conn := ConnectDB(t)
	store := NewStore(conn)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	// Run n concurrent transfer transactions
	n := 10
	amount := int64(10)
	errs := make(chan error)
	fmt.Println("Acc1 and acc2", account1.Balance, account2.Balance)

	// run n concurrent transfer transaction
	for i := 0; i <= n; i++ {

		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 0 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			transResult, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})

			fmt.Printf("fromAccountID: %v ToAccountID: %v\n", transResult.FromAccount.Balance, transResult.ToAccount.Balance)

			errs <- err
		}()

	}

	// Check results
	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err, fmt.Sprint(err))

	}

	// Check the final update balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err, fmt.Sprintf("%v", err))

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	assert.NoError(t, err)

	fmt.Println("After updateAcc, updateacc2 ", updatedAccount1.Balance, updateAccount2.Balance)

	assert.Equal(t, account1.Balance, updatedAccount1.Balance)
	assert.Equal(t, account2.Balance, updateAccount2.Balance)
}
