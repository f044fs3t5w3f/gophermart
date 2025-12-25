package db

import (
	"context"
	"database/sql"

	"github.com/f044fs3t5w3f/gophermart/internal/models"
)

func (d *dbRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	result := d.db.QueryRowContext(ctx, `
		INSERT INTO orders (user_id, number, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (number) DO NOTHING
		RETURNING id`,
		order.UserID, order.Number, order.Status)
	err := result.Err()
	if err != nil {
		return err
	}
	err = result.Scan(&order.ID)
	return err
}

func (d *dbRepository) GetOrderByNumber(ctx context.Context, number string) (*models.Order, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT id, user_id, number, status 
	FROM orders 
	WHERE number = $1`, number)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	order := &models.Order{}
	err = row.Scan(&order.ID, &order.UserID, &order.Number, &order.Status)
	// TODO: /multiplier
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (d *dbRepository) GetOrdersByUserId(ctx context.Context, userId int64) ([]*models.Order, error) {
	rows, err := d.db.QueryContext(ctx, `
	SELECT id, user_id, number, status, uploaded_at, accrual
	FROM orders 
	WHERE user_id = $1
	ORDER BY id DESC`, userId)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*models.Order, 0), nil
		} else {
			return nil, err
		}
	}
	orders := []*models.Order{}
	for rows.Next() {
		rows.Scan()
		order := &models.Order{}
		var accural int32
		err = rows.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.UploadedAt, &accural)
		order.Accural = float64(accural) / multiplier
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (d *dbRepository) ListOrdersForUpdate(ctx context.Context) ([]*models.Order, error) {
	rows, err := d.db.QueryContext(ctx, `
	SELECT id, user_id, number, status, uploaded_at, accrual
	FROM orders 
	WHERE status in ($1, $2)
	ORDER BY id DESC`, models.OrderStatusNew, models.OrderStatusProcessing)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		if err == sql.ErrNoRows {
			return make([]*models.Order, 0), nil
		} else {
			return nil, err
		}
	}
	orders := []*models.Order{}
	for rows.Next() {
		rows.Scan()
		order := &models.Order{}
		var accural int32
		err = rows.Scan(&order.ID, &order.UserID, &order.Number, &order.Status, &order.UploadedAt, &accural)
		order.Accural = float64(accural) / multiplier
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (d *dbRepository) UpdateOrderStatus(ctx context.Context, orderId int64, status models.OrderStatus) error {
	result := d.db.QueryRowContext(ctx, `
		UPDATE orders
		SET status = $1
		WHERE id = $2`,
		status, orderId)
	return result.Err()
}

func (d *dbRepository) UpdateOrderStatusAndAccrual(ctx context.Context, orderId int64, status models.OrderStatus, accrual float64) error {
	accrualInt := int(accrual * multiplier)
	tx, err := d.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})
	defer tx.Rollback()
	if err != nil {
		return err
	}
	result := tx.QueryRowContext(ctx, `
		UPDATE orders
		SET status = $1, accrual = $2
		WHERE id = $3`,
		status, accrualInt, orderId)
	if err = result.Err(); err != nil {
		return err
	}

	var userId int64
	err = tx.QueryRowContext(ctx, `
		SELECT user_id 
		FROM orders 
		WHERE id = $1`, orderId).Scan(&userId)

	if err != nil {
		return err
	}
	var accruals int64
	err = tx.QueryRowContext(ctx, `
		SELECT SUM(accrual) 
		FROM orders 
		WHERE user_id = $1 AND status = $2`, userId, models.OrderStatusProcessed).Scan(&accruals)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `UPDATE users
	SET accruals = $1
	WHERE id = $2
	`, accruals, userId)
	if err != nil {
		return err
	}
	return tx.Commit()
}
