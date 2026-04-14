package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/gfdmit/web-forum/auth-service/config"
	"github.com/gfdmit/web-forum/auth-service/internal/model"
	"github.com/gfdmit/web-forum/auth-service/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresRepository struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, conf config.Postgres) (repository.Repository, error) {
	url := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=disable",
		conf.User, conf.Pass, conf.Host, conf.Port, conf.DB,
	)

	poolConf, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	poolConf.MaxConns = conf.Pool.MaxConns
	poolConf.MinConns = conf.Pool.MinConns
	poolConf.MaxConnLifetime = conf.Pool.MaxConnLifetime
	poolConf.MaxConnIdleTime = conf.Pool.MaxConnIdleTime
	poolConf.HealthCheckPeriod = conf.Pool.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.NewWithConfig: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pool.Ping: %w", err)
	}

	return &postgresRepository{
		db: pool,
	}, nil
}

func (pr *postgresRepository) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	const query = `
        SELECT id, login, password_hash
        FROM forum.users
        WHERE login = $1
    `
	var user model.User
	err := pr.db.QueryRow(ctx, query, login).Scan(
		&user.ID, &user.Login, &user.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, fmt.Errorf("user not found: %w", repository.ErrNotFound)
		}
		return model.User{}, fmt.Errorf("GetUserByLogin: %w", err)
	}
	return user, nil
}

func (pr *postgresRepository) CreateOrUpdateUser(ctx context.Context, input model.CreateUserInput) (model.User, error) {
	const query = `
        INSERT INTO forum.users (login, password_hash)
        VALUES ($1, $2)
        ON CONFLICT (login) DO UPDATE SET password_hash = EXCLUDED.password_hash
        RETURNING id, login, password_hash
    `
	var user model.User
	err := pr.db.QueryRow(ctx, query, input.Login, input.PasswordHash).Scan(
		&user.ID, &user.Login, &user.PasswordHash,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("CreateOrUpdateUser: %w", err)
	}
	return user, nil
}

func (pr *postgresRepository) CreateOrUpdateProfile(ctx context.Context, input model.CreateProfileInput) error {
	const query = `
        INSERT INTO forum.profiles
            (user_id, university_id, firstname, lastname, middlename, birthday, faculty, grade, "group", status)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (user_id) DO UPDATE SET
            university_id = EXCLUDED.university_id,
            firstname     = EXCLUDED.firstname,
            lastname      = EXCLUDED.lastname,
            middlename    = EXCLUDED.middlename,
            birthday      = EXCLUDED.birthday,
            faculty       = EXCLUDED.faculty,
            grade         = EXCLUDED.grade,
            "group"       = EXCLUDED."group",
            status        = EXCLUDED.status,
            updated_at    = NOW()
    `
	_, err := pr.db.Exec(ctx, query,
		input.UserID,
		input.UniversityID,
		input.Firstname,
		input.Lastname,
		input.Middlename,
		input.Birthday,
		input.Faculty,
		input.Grade,
		input.Group,
		input.Status,
	)
	if err != nil {
		return fmt.Errorf("CreateOrUpdateProfile: %w", err)
	}
	return nil
}
