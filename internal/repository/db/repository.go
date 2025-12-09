package db

import (
	"context"
	"database/sql"

	"github.com/f044fs3t5w3f/gophermart/internal/models"
	"github.com/f044fs3t5w3f/gophermart/internal/repository"
)

type dbRepository struct {
	db *sql.DB
}

func (d *dbRepository) CreateSession(ctx context.Context, session *models.Session) error {
	_, err := d.db.ExecContext(ctx, `
		INSERT INTO sessions (user_id, token)
		VALUES ($1, $2)
		ON CONFLICT (token) DO NOTHING`,
		session.UserId, session.Token)
	return err
}

func NewDBRepository(db *sql.DB) repository.Repository {
	return &dbRepository{
		db: db,
	}
}

func (d *dbRepository) CreateUser(ctx context.Context, user *models.User) error {
	result := d.db.QueryRowContext(ctx, `
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2)
		ON CONFLICT (login)  DO NOTHING
		RETURNING id`,
		user.Login, user.PasswordHash)
	err := result.Err()
	if err != nil {
		return err
	}
	err = result.Scan(&user.Id)
	return err
}

func (d *dbRepository) DoesUserExist(ctx context.Context, login string) (bool, error) {
	row := d.db.QueryRowContext(ctx, "SELECT count(*) FROM users where login = $1", login)
	err := row.Err()
	if err != nil {
		return false, err
	}
	var value bool
	err = row.Scan(&value)
	return value, err
}

func (d *dbRepository) GetUserById(context.Context, int64) (*models.User, error) {
	panic("unimplemented")
}

func (d *dbRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	row := d.db.QueryRowContext(ctx, "SELECT id, login, password_hash FROM users where login = $1", login)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	user := &models.User{}
	err = row.Scan(&user.Id, &user.Login, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

func (d *dbRepository) GetUserByToken(ctx context.Context, token string) (*models.User, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT u.id, u.login FROM sessions s
	left JOIN users u
	ON s.user_id = u.id
	WHERE token = $1
`, token)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	user := &models.User{}
	err = row.Scan(&user.Id, &user.Login)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *dbRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	result := d.db.QueryRowContext(ctx, `
		INSERT INTO orders (user_id, number, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (number) DO NOTHING
		RETURNING id`,
		order.UserId, order.Number, order.Status)
	err := result.Err()
	if err != nil {
		return err
	}
	err = result.Scan(&order.Id)
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
	err = row.Scan(&order.Id, &order.UserId, &order.Number, &order.Status)
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
	SELECT id, user_id, number, status, uploaded_at 
	FROM orders 
	WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}
	orders := []*models.Order{}
	for rows.Next() {
		rows.Scan()
		order := &models.Order{}
		err = rows.Scan(&order.Id, &order.UserId, &order.Number, &order.Status, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
