package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.saveUrl"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?,?)")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)

	if err != nil {
		// @TODO свитч-кейс ошибок, например, что такая ссылка не существует
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.sqlite.getUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias =?")

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resUrl string

	err = stmt.QueryRow(alias).Scan(&resUrl)

	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resUrl, nil
}

func (s *Storage) DeleteAlias(alias string) (int64, error) {
	const op = "storage.sqlite.deleteAlias"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias =?")

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(alias)

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return res.RowsAffected()
}
