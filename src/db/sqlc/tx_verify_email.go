package db

import (
	"context"
	"database/sql"
)

// queries signature take an Sql Queries and return an error
// type queries func(*Queries) error

// VerifyEmailrTxParams contains the input parameters of the verify user's email address transaction
type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

// CreateUserTxResult is the result of a successful create user transaction
type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	transactions := func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		result.User, err =  q.UpdateUser(ctx, UpdateUserParams{
			Username: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
		})

		return err
	}

	err := store.execTx(ctx, transactions)

	return result, err
}
