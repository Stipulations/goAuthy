package database

import (
	"context"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

//go:embed schema.sql
var schema string

func New(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func Migrate(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, schema)
	return err
}

func CreateAdmin(ctx context.Context, pool *pgxpool.Pool) error {
	password, err := genPassword()
	if err != nil {
		return  err
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var id uuid.UUID
	err = pool.QueryRow(ctx, `
			INSERT INTO admins (id, username, password_hash)
			VALUES ($1, $2, $3)
			ON CONFLICT (username) DO NOTHING
			RETURNING id`, uuid.New(), "admin", string(hashed),
	).Scan(&id)

	if err == pgx.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Printf("Admin account created. \n Username: admin \n Password: %s \n", password)
	return nil
}


func genPassword() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}