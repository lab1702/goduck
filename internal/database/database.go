package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/marcboeker/go-duckdb"
	"github.com/sirupsen/logrus"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dbPath string, maxConnections int) (*DB, error) {
	// Use connection configuration for read-only access
	dsn := fmt.Sprintf("%s?access_mode=read_only", dbPath)

	conn, err := sql.Open("duckdb", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	conn.SetMaxOpenConns(maxConnections)
	conn.SetMaxIdleConns(maxConnections / 2)
	conn.SetConnMaxLifetime(time.Hour)

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"database_path":   dbPath,
		"max_connections": maxConnections,
	}).Info("Database connection established")

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) GetConnection() *sql.DB {
	return db.conn
}
