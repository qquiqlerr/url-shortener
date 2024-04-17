package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/v5/stdlib"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(DBAddress string) (*Storage, error) {
	const op = "storage.postgres.New"
	db, err := sql.Open("pgx", DBAddress)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS url (
	    id SERIAL PRIMARY KEY,
	    alias TEXT NOT NULL UNIQUE,
	    url TEXT NOT NULL);`)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);`)
	if err != nil {
		return nil, fmt.Errorf("%s - %s", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (id int64, err error) {
	const op = "storage.postgres.SaveURL"
	stmt, err := s.db.Prepare(`INSERT INTO url (alias, url) VALUES ($1, $2) returning id;`)
	if err != nil {
		return 0, fmt.Errorf("%s - %s", op, err)
	}
	err = stmt.QueryRow(alias, urlToSave).Scan(&id)
	if err != nil {
		if pgxErr, ok := err.(pgx.PgError); ok && pgxErr.SQLState() == "23505" {
			return 0, storage.ErrUrlExists
		}
		return 0, fmt.Errorf("execution is failed: %s - %s", op, err)
	}
	return id, nil
}

func (s *Storage) GetURL(alias string) (url string, err error) {
	const op = "storage.postgres.GetURL"

	stmt, err := s.db.Prepare(`SELECT url FROM url WHERE alias = $1;`)
	if err != nil {
		return "", fmt.Errorf("%s - %s", op, err)
	}

	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrUrlNotFound
		}
		return "", fmt.Errorf("%s - %s", op, err)
	}

	return url, nil
}
