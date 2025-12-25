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
	defer tx.Rollback()

	sum := int64(withdraw.Sum * multiplier)

	row := d.db.QueryRowContext(ctx, `
	SELECT accruals - withdraws as balance 
	FROM users
	WHERE id = $1
	FOR UPDATE
`, withdraw.UserId)
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
	var userId int64
	err = row.Scan(&userId)
	if err != sql.ErrNoRows {
		if userId == withdraw.UserId {
			return repository.ErrWithdrawAllreadyExists
		} else {
			return repository.ErrWithdrawAllreadyExistsForAnotherUser
		}
	}
	if err != nil {
		return err
	}

	result := tx.QueryRowContext(ctx, `
		INSERT INTO withdraws (user_id, order_, sum)
		VALUES ($1, $2, $3)
		ON CONFLICT (order_) DO NOTHING
		RETURNING id`,
		withdraw.UserId, withdraw.Order, sum)
	err = result.Err()
	if err != nil {
		return err
	}
	err = result.Scan(&withdraw.Id)
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
	`, userId, userId)
	if err != nil {
		return err
	}

	tx.Commit()
	return err
}

func (d *dbRepository) GetWithdrawsByUserId(ctx context.Context, userId int64) ([]*models.Withdraw, error) {
	rows, err := d.db.QueryContext(ctx, `
	SELECT id, user_id, order_, processed_at, sum
	FROM withdraws 
	WHERE user_id = $1
	ORDER BY id DESC`, userId)
	if err != nil {
		return nil, err
	}
	withdraws := []*models.Withdraw{}
	for rows.Next() {
		rows.Scan()
		withdraw := &models.Withdraw{}
		var sum int64
		err = rows.Scan(&withdraw.Id, &withdraw.UserId, &withdraw.Order, &withdraw.ProcessedAt, &sum)
		withdraw.Sum = float64(sum) / multiplier
		if err != nil {
			return nil, err
		}
		withdraws = append(withdraws, withdraw)
	}
	return withdraws, nil
}
