# GoDuck API Reference

Complete API reference for GoDuck DuckDB REST API Server.

**For getting started and examples, see [GETTING_STARTED.md](GETTING_STARTED.md)**

## Base URL
```
http://localhost:8080
```

## Rate Limiting & Security
- **Rate Limit**: 60 requests per minute per IP address
- **Response**: HTTP 429 when exceeded  
- **Query Size Limit**: Maximum 10KB per SQL query
- **Database Access**: 
  - File databases: Read-only (SELECT statements only)
  - In-memory databases: Read-write (all SQL operations allowed)
- **Content-Type**: All requests/responses use `application/json`

---

## Endpoints

### 1. Execute Query
Execute a SQL query against the DuckDB database. File databases support read-only queries (SELECT), while in-memory databases support all SQL operations.

**URL**: `/query`  
**Method**: `POST`  
**Content-Type**: `application/json`

#### Request Body
```json
{
  "sql": "SELECT column1, column2 FROM table_name WHERE condition LIMIT 100"
}
```

| Field | Type | Required | Description | Limits |
|-------|------|----------|-------------|---------|
| `sql` | string | Yes | SQL query to execute | Max 10KB |

#### Response (Success)
**Status**: `200 OK`
```json
{
  "columns": ["column1", "column2"],
  "rows": [
    ["value1", "value2"],
    ["value3", "value4"]
  ],
  "count": 2,
  "execution_time": "5.2ms"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `columns` | array[string] | Column names in result order |
| `rows` | array[array] | Data rows, each containing values matching columns |
| `count` | integer | Number of rows returned |
| `execution_time` | string | Query execution duration |

#### Response (Error)
**Status**: `400 Bad Request` / `500 Internal Server Error`
```json
{
  "error": "Query execution failed",
  "timestamp": "2025-07-21T19:08:27-04:00"
}
```

#### Example cURL
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{
    "sql": "SELECT name, age FROM users WHERE age > 25 ORDER BY age DESC LIMIT 10"
  }'
```

---

### 2. Health Check
Check if the server and database are healthy and responding.

**URL**: `/health`  
**Method**: `GET`

#### Response (Healthy)
**Status**: `200 OK`
```json
{
  "status": "healthy",
  "timestamp": "2025-07-21T19:08:27-04:00"
}
```

#### Response (Unhealthy)
**Status**: `503 Service Unavailable`
```json
{
  "error": "Database not available",
  "timestamp": "2025-07-21T19:08:27-04:00"
}
```

#### Example cURL
```bash
curl http://localhost:8080/health
```

---

### 3. System Metrics
Get detailed system and database metrics for monitoring.

**URL**: `/metrics`  
**Method**: `GET`

#### Response
**Status**: `200 OK`
```json
{
  "uptime": "2h15m30s",
  "goroutines": 12,
  "memory": {
    "alloc_bytes": 2048576,
    "total_alloc_bytes": 15728640,
    "sys_bytes": 8388608,
    "num_gc": 5
  },
  "database": {
    "max_open_connections": 10,
    "open_connections": 3,
    "in_use": 1,
    "idle": 2
  },
  "timestamp": "2025-07-21T19:08:34-04:00"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `uptime` | string | Server uptime duration |
| `goroutines` | integer | Number of active goroutines |
| `memory.alloc_bytes` | integer | Currently allocated memory |
| `memory.total_alloc_bytes` | integer | Total allocated memory (cumulative) |
| `memory.sys_bytes` | integer | System memory obtained from OS |
| `memory.num_gc` | integer | Number of garbage collection cycles |
| `database.max_open_connections` | integer | Maximum allowed connections |
| `database.open_connections` | integer | Currently open connections |
| `database.in_use` | integer | Connections currently executing queries |
| `database.idle` | integer | Idle connections in pool |

#### Example cURL
```bash
curl http://localhost:8080/metrics
```

---

## Error Codes

| HTTP Status | Description | Common Causes |
|-------------|-------------|---------------|
| `200` | Success | Query executed successfully |
| `400` | Bad Request | Invalid SQL, empty query, query too large |
| `429` | Too Many Requests | Rate limit exceeded |
| `500` | Internal Server Error | Database error, server panic |
| `503` | Service Unavailable | Database not available |

## Error Response Format
All errors return a consistent JSON format:
```json
{
  "error": "Human-readable error description",
  "timestamp": "2025-07-21T19:08:27-04:00"
}
```

## Request Headers
- `Content-Type: application/json` (required for POST requests)
- `X-Request-ID: <uuid>` (optional, for request tracing)

## Response Headers
- `Content-Type: application/json`
- `X-Request-ID: <uuid>` (echoed or generated)
- `Access-Control-Allow-Origin: *` (CORS enabled)

## Query Limitations
- **File Databases**: Read-only access (SELECT statements only)
- **In-Memory Databases**: Full read-write access (all SQL operations)
- **Size Limit**: Maximum 10KB per query
- **Timeout**: Queries timeout after configured duration (default: 30s)
- **No Prepared Statements**: Each request is a single query

## Best Practices
1. **Use LIMIT**: Always limit result sets for better performance
2. **Index Awareness**: Understand your data structure for optimal queries
3. **Error Handling**: Always handle both success and error responses
4. **Rate Limiting**: Implement client-side rate limiting
5. **Monitoring**: Use `/health` and `/metrics` for monitoring
