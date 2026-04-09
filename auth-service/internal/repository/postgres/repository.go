package postgres

import (
	"database/sql"
	"fmt"

	"github.com/gfdmit/web-forum/auth-service/config"

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
