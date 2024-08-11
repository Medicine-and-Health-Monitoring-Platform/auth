package postgres

import (
	"Auth/config"
	"Auth/storage"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func ConnectDB() (storage.IStorage, error) {
	cfg := config.Load()
	conn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_NAME, cfg.DB_PASSWORD)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return Postgres{db: db}, nil
}

func (p Postgres) Close() {
	p.db.Close()
}

func (p Postgres) Admin() storage.IAdminStorage {
	return NewAdminRepo(p.db)
}

func (p Postgres) User() storage.IUserStorage {
	return NewUserRepo(p.db)
}

func (p Postgres) Token() storage.ITokenStorage {
	return NewTokenRepo(p.db)
}
