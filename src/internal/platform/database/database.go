package database

import (
	"context"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
)

type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

func Open(cfg Config) (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("timezone", "utc")

	if cfg.DisableTLS {
		q.Set("sslmode", "disable")
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}
	return sqlx.Open("postgres", u.String())
}

func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
