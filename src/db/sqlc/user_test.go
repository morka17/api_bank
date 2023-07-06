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




func TestUpdateUserOnlyFullName(t *testing.T){
	oldUser := CreateRandomUser(t)
	
	newFullName := utils.RandomOwner()
	updated, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid: true,
		},
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, updated.FullName, newFullName)
	assert.Equal(t, newFullName, updated.FullName)
}


func TestUpdateUserOnlyEmail(t *testing.T){
	oldUser := CreateRandomUser(t)
	
	email := utils.RandomEmail()
	updated, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: email,
			Valid: true,
		},
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, updated.Email, oldUser.Email)
	assert.Equal(t, email, updated.Email)
}

func TestUpdateUserOnlyHashPassword(t *testing.T){
	oldUser := CreateRandomUser(t)
	
	hashPassword, _ := utils.HashPassword("12345676")
	updated, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: hashPassword,
			Valid: true,
		},
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, updated.HashedPassword, oldUser.HashedPassword)
	assert.Equal(t, hashPassword, updated.HashedPassword)
}



func TestUpdateUserOnlyAllField(t *testing.T){
	oldUser := CreateRandomUser(t)
	
	hashPassword, _ := utils.HashPassword("12345676")
	newFull := utils.RandomOwner()
	email := utils.RandomEmail()
	updated, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: hashPassword,
			Valid: true,
		},
		Email: sql.NullString{
			String: email,
			Valid: true,
		},
		FullName: sql.NullString{
			String: newFull,
			Valid: true,
		},
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, updated.HashedPassword, oldUser.HashedPassword)
	assert.Equal(t, hashPassword, updated.HashedPassword)
	assert.NotEmpty(t, updated.Email, oldUser.Email)
	assert.Equal(t, email, updated.Email)
	assert.NotEmpty(t, updated.FullName, oldUser.FullName)
	assert.Equal(t, newFull, updated.FullName)
}