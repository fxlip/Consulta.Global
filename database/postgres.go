package database

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Função para buscar a senha do banco
func getDBPassword() string {
	if path := os.Getenv("PG_PASSWORD_FILE"); path != "" {
		data, err := os.ReadFile(path)
		if err == nil {
			return strings.TrimSpace(string(data))
		}
	}
	return os.Getenv("PG_PASSWORD")
}

func ConnectPostgres() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("PG_USER"),
		getDBPassword(),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_DBNAME"),
	))

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	return pool, err
}
