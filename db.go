package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func InitDB(cfg DatabaseConfig) error {
	config, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return fmt.Errorf("unable to parse database config: %w", err)
	}

	config.MaxConns = int32(cfg.MaxConns)

	db, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Connected to database")
	return nil
}

func LogBotAction(level, message, details string) {
	if db == nil {
		log.Println("Database not initialized, skipping log:", message)
		return
	}
	_, err := db.Exec(context.Background(), "INSERT INTO bot_logs (level, message, details) VALUES ($1, $2, $3)", level, message, details)
	if err != nil {
		log.Println("Failed to write log to DB:", err)
	}
}
