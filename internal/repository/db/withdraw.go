package db

import (
	"context"
	"database/sql"

	"github.com/f044fs3t5w3f/gophermart/internal/models"
	"github.com/f044fs3t5w3f/gophermart/internal/repository"
)

func (d *dbRepository) CreateWithdraw(ctx context.Context, withdraw *models.Withdraw) error {
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	sum := int64(withdraw.Sum * multiplier)

	row := tx.QueryRowContext(ctx, `
	SELECT accruals - withdraws as balance 
	FROM users
	WHERE id = $1
	FOR UPDATE
`, withdraw.UserID)
	var balance int64
	err = row.Scan(&balance)
	if err != nil {
		return err
	}
	if balance < sum {
		return repository.ErrNotEnough
	}

	row = tx.QueryRowContext(ctx, `
		SELECT user_id
		FROM withdraws
		WHERE order_ = $1
	`, withdraw.Order)
	var userID int64
	err = row.Scan(&userID)

	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	} else {
		if userID == withdraw.UserID {
			return repository.ErrWithdrawAllreadyExists
		} else {
			return repository.ErrWithdrawAllreadyExistsForAnotherUser
		}
	}

	result := tx.QueryRowContext(ctx, `
		INSERT INTO withdraws (user_id, order_, sum)
		VALUES ($1, $2, $3)
		ON CONFLICT (order_) DO NOTHING
		RETURNING id`,
		withdraw.UserID, withdraw.Order, sum)
	err = result.Err()
	if err != nil {
		return err
	}
	err = result.Scan(&withdraw.ID)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE users
	SET withdraws = (
		SELECT SUM(sum)
		FROM withdraws
		WHERE user_id = $1
	)
	WHERE id = $2
	`, withdraw.UserID, withdraw.UserID)
	if err != nil {
		return err
	}

	tx.Commit()
	return err
}

func (d *dbRepository) GetWithdrawsByUserID(ctx context.Context, userID int64) ([]*models.Withdraw, error) {
	rows, err := d.db.QueryContext(ctx, `
	SELECT id, user_id, order_, processed_at, sum
	FROM withdraws 
	WHERE user_id = $1
	ORDER BY id DESC`, userID)
	if err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		if err != sql.ErrNoRows {
			return nil, nil
		}
	}
	withdraws := []*models.Withdraw{}
	for rows.Next() {
		rows.Scan()
		withdraw := &models.Withdraw{}
		var sum int64
		err = rows.Scan(&withdraw.ID, &withdraw.UserID, &withdraw.Order, &withdraw.ProcessedAt, &sum)
		withdraw.Sum = float64(sum) / multiplier
		if err != nil {
			return nil, err
		}
		withdraws = append(withdraws, withdraw)
	}
	return withdraws, nil
}
