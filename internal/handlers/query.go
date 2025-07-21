package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"goduck/internal/database"
	"goduck/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type QueryHandler struct {
	db           *database.DB
	queryTimeout time.Duration
}

func NewQueryHandler(db *database.DB, timeout time.Duration) *QueryHandler {
	return &QueryHandler{
		db:           db,
		queryTimeout: timeout,
	}
}

func (h *QueryHandler) ExecuteQuery(c *gin.Context) {
	var req models.QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: fmt.Sprintf("Invalid request: %v", err),
			Time:  time.Now(),
		})
		return
	}

	if strings.TrimSpace(req.SQL) == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "SQL query cannot be empty",
			Time:  time.Now(),
		})
		return
	}

	// Limit query size to prevent abuse
	if len(req.SQL) > 10000 { // 10KB limit
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "SQL query too large (max 10KB)",
			Time:  time.Now(),
		})
		return
	}

	start := time.Now()

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.queryTimeout)
	defer cancel()

	rows, err := h.db.GetConnection().QueryContext(ctx, req.SQL)
	if err != nil {
		requestID, _ := c.Get("request_id")
		logrus.WithFields(logrus.Fields{
			"request_id": requestID,
			"sql":        req.SQL,
			"error":      err.Error(),
		}).Error("Query execution failed")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Query execution failed",
			Time:  time.Now(),
		})
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		logrus.WithError(err).Error("Failed to get columns")
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to process query results",
			Time:  time.Now(),
		})
		return
	}

	var result [][]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			logrus.WithError(err).Error("Failed to scan row")
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error: "Failed to process query results",
				Time:  time.Now(),
			})
			return
		}

		for i, val := range values {
			if val == nil {
				values[i] = nil
			} else if b, ok := val.([]byte); ok {
				values[i] = string(b)
			}
		}

		result = append(result, values)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: fmt.Sprintf("Row iteration error: %v", err),
			Time:  time.Now(),
		})
		return
	}

	duration := time.Since(start)

	requestID, _ := c.Get("request_id")
	logrus.WithFields(logrus.Fields{
		"request_id":     requestID,
		"sql":            req.SQL,
		"execution_time": duration,
		"row_count":      len(result),
	}).Info("Query executed successfully")

	c.JSON(http.StatusOK, models.QueryResponse{
		Columns: columns,
		Rows:    result,
		Count:   len(result),
		Time:    duration.String(),
	})
}

func (h *QueryHandler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := h.db.GetConnection().PingContext(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, models.ErrorResponse{
			Error: fmt.Sprintf("Database not available: %v", err),
			Time:  time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, models.HealthResponse{
		Status: "healthy",
		Time:   time.Now().Format(time.RFC3339),
	})
}
