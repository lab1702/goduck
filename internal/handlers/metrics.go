package handlers

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type MetricsResponse struct {
	Uptime     string        `json:"uptime"`
	Goroutines int           `json:"goroutines"`
	Memory     MemoryStats   `json:"memory"`
	Database   DatabaseStats `json:"database"`
	Timestamp  string        `json:"timestamp"`
}

type MemoryStats struct {
	Alloc      uint64 `json:"alloc_bytes"`
	TotalAlloc uint64 `json:"total_alloc_bytes"`
	Sys        uint64 `json:"sys_bytes"`
	NumGC      uint32 `json:"num_gc"`
}

type DatabaseStats struct {
	MaxOpenConnections int `json:"max_open_connections"`
	OpenConnections    int `json:"open_connections"`
	InUse              int `json:"in_use"`
	Idle               int `json:"idle"`
}

var startTime = time.Now()

func (h *QueryHandler) Metrics(c *gin.Context) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	dbStats := h.db.GetConnection().Stats()

	metrics := MetricsResponse{
		Uptime:     time.Since(startTime).String(),
		Goroutines: runtime.NumGoroutine(),
		Memory: MemoryStats{
			Alloc:      m.Alloc,
			TotalAlloc: m.TotalAlloc,
			Sys:        m.Sys,
			NumGC:      m.NumGC,
		},
		Database: DatabaseStats{
			MaxOpenConnections: dbStats.MaxOpenConnections,
			OpenConnections:    dbStats.OpenConnections,
			InUse:              dbStats.InUse,
			Idle:               dbStats.Idle,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, metrics)
}
