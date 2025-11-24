package postgres

import (
	"context"
	"log"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

    "github.com/kgugunava/avito-tech-internship/internal/config"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func NewPostgres() Postgres{
	return Postgres{
		Pool: &pgxpool.Pool{},
	}
}

func (p *Postgres) ConnectToPostgresMainDatabase(cfg config.Config) error { // для подключения к бд постгреса для создания нашей бд
    dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
    cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, "postgres", cfg.SslMode)
    
    fmt.Println("Connecting to DB with:", dbUrl)
    
    newPostgresPool, err := pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        log.Fatal("Eror while connecting to postgres database\n", err)
        return err
    }
    p.Pool = newPostgresPool
    return nil
}

func (p *Postgres) ConnectToDatabase(cfg config.Config) error { // для подключения к нужной бд
	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
    cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.SslMode)
    
    newPostgresPool, err := pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        log.Fatal("Error while connecting to database\n", err)
        return err
    }
    p.Pool = newPostgresPool
    return nil
} 

func (p *Postgres) CreateDatabase(cfg config.Config) error {
	var dbExists bool
    err := p.Pool.QueryRow(context.Background(), 
        "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.DbName).Scan(&dbExists)
    if err != nil {
        log.Fatal(err)
    }
    
    if !dbExists {
        _, err := p.Pool.Exec(context.Background(), 
            fmt.Sprintf("CREATE DATABASE %s", cfg.DbName))
        if err != nil {
            log.Fatal(err)
            return err
        }

        p.Pool.Close()
    
        dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", 
            cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName, cfg.SslMode)
        
        p.Pool, err = pgxpool.New(context.Background(), dbUrl)
        if err != nil {
            log.Fatal("Eror while creating database\n", err)
            return err
        }
    }
    return nil
}

func (p *Postgres) CreateDatabaseTables() error {
    _, err := p.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS teams (
            team_id SERIAL PRIMARY KEY,
            team_name TEXT UNIQUE NOT NULL
        );
    `)
	if err != nil {
		log.Fatal("Error creating teams table\n", err)
		return err
	}

	_, err = p.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS users (
            user_id TEXT PRIMARY KEY,
            username TEXT NOT NULL,
            team_id INT REFERENCES teams(team_id) ON DELETE SET NULL,
            is_active BOOLEAN NOT NULL
        );
    `)
	if err != nil {
		log.Fatal("Error creating users table\n", err)
		return err
	}

	_, err = p.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS pull_requests (
            pull_request_id TEXT PRIMARY KEY,
            pull_request_name TEXT NOT NULL,
            author_id TEXT REFERENCES users(user_id),
            status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')),
            created_at TIMESTAMPTZ DEFAULT NOW(),
            merged_at TIMESTAMPTZ
        );
    `)
	if err != nil {
		log.Fatal("Error creating pull_requests table\n", err)
		return err
	}

	_, err = p.Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS pull_request_reviewers (
            pull_request_id TEXT REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
            reviewer_id TEXT REFERENCES users(user_id) ON DELETE CASCADE,
            PRIMARY KEY (pull_request_id, reviewer_id)
        );
    `)
	if err != nil {
		log.Fatal("Error creating pull_request_reviewers table\n", err)
		return err
	}

	return nil
}
