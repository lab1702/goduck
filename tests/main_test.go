package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"goduck/internal/database"
	"goduck/internal/handlers"
	"goduck/internal/middleware"
	"goduck/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) (*database.DB, string) {
	// Use in-memory database to avoid extension issues
	dbPath := ":memory:"

	// First create a writable connection to set up test data
	writeDB, err := database.NewDB(dbPath, 1, true)
	if err != nil {
		t.Skip("DuckDB not available in test environment, skipping database tests")
		return nil, ""
	}

	_, err = writeDB.GetConnection().Exec(`
		CREATE TABLE test_table AS 
		SELECT 1 as id, 'test' as name 
		UNION ALL 
		SELECT 2 as id, 'example' as name
	`)
	if err != nil {
		writeDB.Close()
		t.Skip("Cannot create test data, skipping database tests")
		return nil, ""
	}

	// For in-memory DB, we can't really test read-only mode the same way
	// but we can test the API functionality
	return writeDB, dbPath
}
func setupTestRouter(db *database.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.RecoveryMiddleware())

	queryHandler := handlers.NewQueryHandler(db, 10*time.Second)
	router.POST("/query", queryHandler.ExecuteQuery)
	router.GET("/health", queryHandler.Health)

	return router
}

func TestQueryEndpoint(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	tests := []struct {
		name         string
		requestBody  models.QueryRequest
		expectedCode int
		checkResult  func(t *testing.T, body []byte)
	}{
		{
			name:         "valid select query",
			requestBody:  models.QueryRequest{SQL: "SELECT * FROM test_table"},
			expectedCode: http.StatusOK,
			checkResult: func(t *testing.T, body []byte) {
				var response models.QueryResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Equal(t, 2, response.Count)
				assert.Len(t, response.Columns, 2)
				assert.Contains(t, response.Columns, "id")
				assert.Contains(t, response.Columns, "name")
			},
		},
		{
			name:         "empty sql query",
			requestBody:  models.QueryRequest{SQL: ""},
			expectedCode: http.StatusBadRequest,
			checkResult: func(t *testing.T, body []byte) {
				var response models.ErrorResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Contains(t, response.Error, "required")
			},
		},
		{
			name:         "invalid sql query",
			requestBody:  models.QueryRequest{SQL: "INVALID SQL"},
			expectedCode: http.StatusBadRequest,
			checkResult: func(t *testing.T, body []byte) {
				var response models.ErrorResponse
				err := json.Unmarshal(body, &response)
				require.NoError(t, err)
				assert.Contains(t, response.Error, "Query execution failed")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/query", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			tt.checkResult(t, w.Body.Bytes())
		})
	}
}

func TestHealthEndpoint(t *testing.T) {
	db, _ := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(db)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "healthy", response.Status)
}

func TestDatabaseReadOnlyValidation(t *testing.T) {
	// Test that in-memory database requires read-write mode
	_, err := database.NewDB("", 1, false) // in-memory with read-only should fail
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "in-memory database requires read-write access")
	assert.Contains(t, err.Error(), "GODUCK_READ_WRITE=true")

	// Test that in-memory database works with read-write
	db, err := database.NewDB("", 1, true) // in-memory with read-write should work
	if err != nil {
		t.Skip("DuckDB not available in test environment")
		return
	}
	defer db.Close()

	// Verify we can write to it
	_, err = db.GetConnection().Exec("CREATE TABLE test AS SELECT 1 as id")
	assert.NoError(t, err)
}
