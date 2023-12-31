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

var (
	testQueries *Queries
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

func CreateRandomAccount(t *testing.T) Account {

	user := CreateRandomUser(t)


	arg := CreateAccountParams{
		Owner:   user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account)

	assert.Equal(t, arg.Owner, account.Owner)
	assert.Equal(t, arg.Balance, account.Balance)
	assert.Equal(t, arg.Currency, account.Currency)

	assert.NotZero(t, account.ID)
	assert.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {

	CreateRandomAccount(t)
}


func TestGetAccount(t *testing.T){
	account1 := CreateRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, account2)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, account1.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.WithinDuration(t, account1.CreatedAt, account1.CreatedAt, time.Second)
}



func TestUpdateAccount(t *testing.T){

	account1 := CreateRandomAccount(t)
	
	arg := UpdateAccountsParams {
		ID: account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccounts(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, account2)

	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account1.Owner, account2.Owner)
	assert.Equal(t, arg.Balance, account2.Balance)
	assert.Equal(t, account1.Currency, account2.Currency)
	assert.WithinDuration(t, account1.CreatedAt, account1.CreatedAt, time.Second)
}


func TestDeleteAccount(t *testing.T){
	account1 := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	assert.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.Error(t, err)
	assert.Empty(t, account2)
}


func TestListAccounts(t *testing.T){

	var lastAccount Account 

	for i := 0; i < 10; i++ {
		lastAccount = CreateRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner: lastAccount.Owner,
		Limit: 5,
		Offset: 5,
	}
	
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, accounts)

	for _, account := range accounts {
		assert.NotEmpty(t, account)
	}

}