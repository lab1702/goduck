package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/marcboeker/go-duckdb/v2"
	"github.com/sirupsen/logrus"
)

type DB struct {
	conn *sql.DB
}

func NewDB(dbPath string, maxConnections int) (*DB, error) {
	var dsn string
	var inMemory bool

	if dbPath == "" {
		// Use in-memory database with read-write access when no path is specified
		dsn = ":memory:?access_mode=read_write"
		inMemory = true
	} else {
		// Use connection configuration for read-only access for file databases
		dsn = fmt.Sprintf("%s?access_mode=read_only", dbPath)
		inMemory = false
	}

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

	if inMemory {
		logrus.WithFields(logrus.Fields{
			"database_mode":   "in-memory",
			"access_mode":     "read_write",
			"max_connections": maxConnections,
		}).Info("In-memory database connection established")
	} else {
		logrus.WithFields(logrus.Fields{
			"database_path":   dbPath,
			"access_mode":     "read_only",
			"max_connections": maxConnections,
		}).Info("Database connection established")
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) GetConnection() *sql.DB {
	return db.conn
}
