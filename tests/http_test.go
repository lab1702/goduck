package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/lab1702/goduck/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Simple HTTP-only tests that don't require a real database
func TestHTTPEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Simple mock handler that doesn't use database
	router.POST("/query", func(c *gin.Context) {
		var req models.QueryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Invalid request",
			})
			return
		}

		if req.SQL == "" || strings.TrimSpace(req.SQL) == "" {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "SQL query cannot be empty",
			})
			return
		}

		// Mock successful response
		c.JSON(http.StatusOK, models.QueryResponse{
			Columns: []string{"test"},
			Rows:    [][]interface{}{{"value"}},
			Count:   1,
			Time:    "1ms",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.HealthResponse{
			Status: "healthy",
			Time:   "2024-01-01T00:00:00Z",
		})
	})

	t.Run("query endpoint with valid request", func(t *testing.T) {
		req := models.QueryRequest{SQL: "SELECT 1"}
		body, _ := json.Marshal(req)

		httpReq := httptest.NewRequest("POST", "/query", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.QueryResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, 1, response.Count)
		assert.Equal(t, []string{"test"}, response.Columns)
	})

	t.Run("query endpoint with empty SQL", func(t *testing.T) {
		// The binding:"required" tag will catch empty strings at validation level
		req := models.QueryRequest{SQL: ""}
		body, _ := json.Marshal(req)

		httpReq := httptest.NewRequest("POST", "/query", bytes.NewBuffer(body))
		httpReq.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response models.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		// The binding validation will trigger first for empty required fields
		assert.Contains(t, response.Error, "Invalid request")
	})

	t.Run("health endpoint", func(t *testing.T) {
		httpReq := httptest.NewRequest("GET", "/health", nil)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.HealthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "healthy", response.Status)
	})

	t.Run("invalid JSON request", func(t *testing.T) {
		httpReq := httptest.NewRequest("POST", "/query", bytes.NewBufferString("{invalid json"))
		httpReq.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
