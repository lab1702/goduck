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

func NewDB(dbPath string, maxConnections int, readWrite bool) (*DB, error) {
	// Validate configuration
	if dbPath == "" && !readWrite {
		return nil, fmt.Errorf("in-memory database requires read-write access (set GODUCK_READ_WRITE=true)")
	}

	var dsn string
	var accessMode string

	if readWrite {
		accessMode = "read_write"
	} else {
		accessMode = "read_only"
	}

	if dbPath == "" {
		// Use in-memory database
		dsn = fmt.Sprintf(":memory:?access_mode=%s", accessMode)
		dbPath = ":memory:"
	} else {
		// Use file database
		dsn = fmt.Sprintf("%s?access_mode=%s", dbPath, accessMode)
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

	logrus.WithFields(logrus.Fields{
		"database_path":   dbPath,
		"access_mode":     accessMode,
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
