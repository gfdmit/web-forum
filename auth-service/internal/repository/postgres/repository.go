package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gfdmit/web-forum/auth-service/config"
	"github.com/gfdmit/web-forum/auth-service/internal/model"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func New(conf config.Postgres) (*Repository, error) {
	url := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=disable", conf.User, conf.Pass, conf.Host, conf.Port, conf.DB)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db.Ping: %v", err)
	}
	return &Repository{
		db: db,
	}, nil
}

func (r Repository) GetIdPassHash(login string) (int, string, error) {
	var (
		id       int
		passHash string
	)
	err := r.db.QueryRow(
		"SELECT id, password_hash FROM forum.users WHERE login = $1",
		login,
	).Scan(&id, &passHash)

	if errors.Is(err, sql.ErrNoRows) {
		return -1, "", ErrNotFound
	}
	if err != nil {
		return -1, "", err
	}
	return id, passHash, nil
}

func (r Repository) CreateOrUpdateUser(login, newPassHash string) (int, error) {
	var (
		id int
	)
	err := r.db.QueryRow(
		`INSERT INTO forum.users (login, password_hash)
		VALUES($1, $2) 
		ON CONFLICT (login) DO UPDATE SET password_hash = EXCLUDED.password_hash 
		RETURNING id`,
		login, newPassHash,
	).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r Repository) CreateOrUpdateProfile(id int, profile *model.Profile) error {
	_, err := r.db.Exec(
		`INSERT INTO forum.profiles 
            (user_id, university_id, firstname, lastname, middlename, birthday, faculty, grade, "group", status)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
         ON CONFLICT (user_id) DO UPDATE SET
            university_id = EXCLUDED.university_id,
            firstname  = EXCLUDED.firstname,
            lastname   = EXCLUDED.lastname,
            middlename = EXCLUDED.middlename,
            birthday   = EXCLUDED.birthday,
            faculty    = EXCLUDED.faculty,
            grade      = EXCLUDED.grade,
            "group"    = EXCLUDED."group",
            status     = EXCLUDED.status,
            updated_at = NOW()`,
		id,
		profile.AllId,
		profile.Firstname,
		profile.Lastname,
		profile.MiddleName,
		profile.Birthday,
		profile.Faculty,
		profile.Grade,
		profile.Group,
		profile.Status,
	)
	return err
}
