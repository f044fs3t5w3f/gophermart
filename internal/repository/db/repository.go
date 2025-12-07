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
	result, err := d.db.ExecContext(ctx, `
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2)
		ON CONFLICT (login)  DO NOTHING`,
		user.Login, user.PasswordHash)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	user.Id = id
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
	return user, err
}
