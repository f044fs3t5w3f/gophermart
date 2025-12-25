package db

import (
	"context"
	"database/sql"

	"github.com/f044fs3t5w3f/gophermart/internal/models"
)

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
	err = result.Scan(&user.ID)
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

func (d *dbRepository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	row := d.db.QueryRowContext(ctx, "SELECT id, login, password_hash FROM users where login = $1", login)
	err := row.Err()
	if err != nil {
		return nil, err
	}
	user := &models.User{}
	err = row.Scan(&user.ID, &user.Login, &user.PasswordHash)
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
	err = row.Scan(&user.ID, &user.Login)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *dbRepository) GetBalanceByUserId(ctx context.Context, userId int64) (float64, float64, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT accruals, withdraws 
	FROM users
	WHERE id = $1
`, userId)
	var accrualsInt, withdrawsInt int64
	err := row.Scan(&accrualsInt, &withdrawsInt)
	return float64(accrualsInt) / multiplier, float64(withdrawsInt) / multiplier, err
}
