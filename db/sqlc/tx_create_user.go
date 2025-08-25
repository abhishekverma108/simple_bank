package db

import (
	"context"
	"go.elastic.co/apm/v2"
)

 type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error // AfterCreate is a callback function that will be called after creating the user
}

 type CreateUserTxResult struct {
	User User
}

 func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		spanCreate, _ := apm.StartSpan(ctx, "db.create_user", "db")
		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		spanCreate.End()
		if err != nil {
			apm.CaptureError(ctx, err).Send()
			return err
		}
		spanAfter, _ := apm.StartSpan(ctx, "after_create_callback", "custom")
		err = arg.AfterCreate(result.User)
		spanAfter.End()
		if err != nil {
			apm.CaptureError(ctx, err).Send()
		}
		return err
	})
	return result, err
}
