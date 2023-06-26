package db

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/morka17/shiny_bank/v1/src/utils"
	"github.com/stretchr/testify/assert"
)

func init() {
	config, err := utils.LoadConfig("../../..")
	if err != nil {
		log.Fatalf("Failed to load config %v", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("Expected no error, but found %v", err)
	}

	testQueries = New(conn)

}

func CreateRandomUser(t *testing.T) User {
	
	hashedPassword, err := utils.HashPassword(utils.RandomString(6))
	assert.NoError(t, err)

	arg := CreateUserParams{
		Username:          utils.RandomOwner(),
		FullName:          utils.RandomOwner(),
		Email:           utils.RandomEmail(),
		HashedPassword:  hashedPassword,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, user)

	assert.Equal(t, arg.Username, user.Username)
	assert.Equal(t, arg.HashedPassword, user.HashedPassword)
	assert.Equal(t, arg.FullName, user.FullName)
	assert.Equal(t, arg.Email, user.Email)

	assert.True(t, user.PasswordChangeAt.IsZero())
	assert.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {

	CreateRandomAccount(t)
}

func TestGetUser(t *testing.T) {
	user := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user.Username)
	assert.NoError(t, err)
	assert.NotEmpty(t, user2)

	assert.Equal(t, user2.Username, user.Username)
	assert.Equal(t, user2.HashedPassword, user.HashedPassword)
	assert.Equal(t, user2.FullName, user.FullName)
	assert.Equal(t, user2.Email, user.Email)

	assert.WithinDuration(t, user2.CreatedAt, user.CreatedAt, time.Second)
}

// func TestUpdateAccount(t *testing.T){

// 	account1 := CreateRandomAccount(t)

// 	arg := UpdateAccountsParams {
// 		ID: account1.ID,
// 		Balance: utils.RandomMoney(),
// 	}

// 	account2, err := testQueries.UpdateAccounts(context.Background(), arg)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, account2)

// 	assert.Equal(t, account1.ID, account2.ID)
// 	assert.Equal(t, account1.Owner, account2.Owner)
// 	assert.Equal(t, arg.Balance, account2.Balance)
// 	assert.Equal(t, account1.Currency, account2.Currency)
// 	assert.WithinDuration(t, account1.CreatedAt, account1.CreatedAt, time.Second)
// }

// func TestDeleteAccount(t *testing.T){
// 	account1 := CreateRandomAccount(t)
// 	err := testQueries.DeleteAccount(context.Background(), account1.ID)
// 	assert.NoError(t, err)

// 	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
// 	assert.Error(t, err)
// 	assert.Empty(t, account2)
// }

// func TestListAccounts(t *testing.T){

// 	for i := 0; i < 10; i++ {
// 		CreateRandomAccount(t)
// 	}

// 	arg := ListAccountsParams{
// 		Limit: 5,
// 		Offset: 5,
// 	}

// 	accounts, err := testQueries.ListAccounts(context.Background(), arg)
// 	assert.NoError(t, err)
// 	assert.Len(t, accounts, 5)

// 	for _, account := range accounts {
// 		assert.NotEmpty(t, account)
// 	}

// }
