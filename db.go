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

        // Create tables if they don't exist
        queries := []string{
                `CREATE TABLE IF NOT EXISTS snt_users (
                        created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        user_id BIGINT NOT NULL PRIMARY KEY,
                        user_name VARCHAR(64) NOT NULL,
                        user_fio VARCHAR(255),
                        user_phone VARCHAR(10),
                        comment TEXT
                )`,
                `CREATE INDEX IF NOT EXISTS idx_snt_users_user_name ON snt_users(user_name)`,
                `CREATE TABLE IF NOT EXISTS snt_details (
                        created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        id VARCHAR(8) NOT NULL PRIMARY KEY,
                        name VARCHAR(120) NOT NULL,
                        inn VARCHAR(10) NOT NULL,
                        kpp VARCHAR(9) NOT NULL,
                        personal_acc VARCHAR(20) NOT NULL,
                        bank_name VARCHAR(120) NOT NULL,
                        bik VARCHAR(9) NOT NULL,
                        corresp_acc VARCHAR(20) NOT NULL,
                        comment TEXT
                )`,
                `CREATE TABLE IF NOT EXISTS snt_contacts (
                        created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        modified TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        prior INTEGER NOT NULL,
                        type VARCHAR(20) NOT NULL PRIMARY KEY,
                        value VARCHAR(120) NOT NULL,
                        adds VARCHAR(240),
                        comment TEXT
                )`,
                `CREATE TABLE IF NOT EXISTS bot_logs (
                        id SERIAL PRIMARY KEY,
                        created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
                        level VARCHAR(10),
                        message TEXT,
                        details TEXT
                )`,
        }

        for _, q := range queries {
                if _, err := db.Exec(context.Background(), q); err != nil {
                        log.Printf("Warning: failed to execute initialization query: %v", err)
                }
        }

        log.Println("Connected to database and initialized tables")
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
