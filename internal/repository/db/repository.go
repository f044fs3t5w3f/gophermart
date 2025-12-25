package db

import (
	"context"
	"database/sql"

	"github.com/f044fs3t5w3f/gophermart/internal/models"
)

const multiplier = 100000

// type ContextTransactionKeyType string

// var ContextTransactionKey ContextTransactionKeyType = "tx"

type dbRepository struct {
	db *sql.DB
}

// func (d *dbRepository) GetTransactionContext(ctx context.Context) (context.Context, error) {
// 	tx, err := d.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return context.WithValue(ctx, ContextTransactionKey, tx), nil
// }

// func (d *dbRepository) CommitTransaction(ctx context.Context) error {
// 	tx, ok := ctx.Value(ContextTransactionKey).(*sql.Tx)
// 	if !ok {
// 		return errors.New("commit failed")
// 	}
// 	return tx.Commit()
// }

func (d *dbRepository) CreateSession(ctx context.Context, session *models.Session) error {
	_, err := d.db.ExecContext(ctx, `
		INSERT INTO sessions (user_id, token)
		VALUES ($1, $2)
		ON CONFLICT (token) DO NOTHING`,
		session.UserId, session.Token)
	return err
}

func NewDBRepository(db *sql.DB) *dbRepository {
	return &dbRepository{
		db: db,
	}
}
