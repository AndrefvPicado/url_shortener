package repository

import (
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type PostgresRepo struct {
	db *sql.DB
}

var PostgresRepoInstance *PostgresRepo

// NewPostgresRepo initializes a new PostgresRepo with an open database connection.
func NewPostgresRepo(connStr string) (*PostgresRepo, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Verify the connection is established
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepo{db: db}, nil
}

// SaveURL saves the short code and original URL to the database.
func (p *PostgresRepo) SaveURL(shortCode string, originalURL string) error {
	query := `INSERT INTO urls (short_code, original_url) 
              VALUES ($1, $2) 
              ON CONFLICT (short_code) DO NOTHING`

	_, err := p.db.Exec(query, shortCode, originalURL)
	if err != nil {
		return err
	}
	return nil
}

// GetOriginalURL retrieves the original URL using the short code.
func (p *PostgresRepo) GetOriginalURL(shortCode string) (string, error) {
	var originalURL string
	query := `SELECT original_url FROM urls WHERE short_code = $1`
	err := p.db.QueryRow(query, shortCode).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return originalURL, nil
}
