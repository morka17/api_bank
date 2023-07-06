package db

import "context"

// queries signature take an Sql Queries and return an error
// type queries func(*Queries) error

// CreateUserTxParams contains the input parameters of the create user transaction
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// CreateUserTxResult is the result of a successful create user transaction
type CreateUserTxResult struct {
	User User
}

// TransferTx performs a money transfer from one account to the other
// It create a transfer record, add accounts entries, and update accounts' balance within a signle datanase transaction
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	transactions := func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)

	}

	err := store.execTx(ctx, transactions)

	return result, err
}
