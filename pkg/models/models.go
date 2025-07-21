package models

import "time"

type QueryRequest struct {
	SQL string `json:"sql" binding:"required"`
}

type QueryResponse struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
	Count   int             `json:"count"`
	Time    string          `json:"execution_time"`
}

type ErrorResponse struct {
	Error string    `json:"error"`
	Time  time.Time `json:"timestamp"`
}

type HealthResponse struct {
	Status string `json:"status"`
	Time   string `json:"timestamp"`
}
