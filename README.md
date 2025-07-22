# GoDuck - DuckDB REST API Server

A production-ready, high-performance REST API server for executing SQL queries against DuckDB databases. Supports both file-based databases (read-only) and in-memory databases (read-write).

## 📚 Documentation

- 🚀 **New to GoDuck?** Start with [GETTING_STARTED.md](GETTING_STARTED.md)
- 📋 **API Reference** [API.md](API.md)  
- 📝 **Version History** [CHANGELOG.md](CHANGELOG.md)
- 📚 **Documentation Guide** [DOCS.md](DOCS.md)

## ✨ Features

- **🔒 Security First**: Read-only file databases or read-write in-memory databases with rate limiting and input validation
- **⚡ High Performance**: Connection pooling, query timeouts, and optimized execution
- **📊 Production Ready**: Comprehensive monitoring, metrics, and structured logging
- **🛡️ Robust**: Request tracing, graceful shutdown, and error handling
- **🐳 Container Ready**: Docker support with multi-stage builds
- **📈 Observable**: Health checks, metrics endpoint, and request tracing

## 🚀 Quick Start

### Option 1: Use with Your DuckDB File
```bash
# Download and run GoDuck with your database
export GODUCK_DATABASE_PATH="/path/to/your/database.duckdb"
./goduck

# Test it works
curl http://localhost:8080/health
```

### Option 2: Try with In-Memory Database
```bash
# Run with temporary in-memory database (for testing)
export GODUCK_READ_WRITE=true
./goduck

# Test it works
curl http://localhost:8080/health
```

### Option 3: Use Sample Data
```bash
# Create sample database and run
./start_example.sh
export GODUCK_DATABASE_PATH="./data/sample.duckdb"
./goduck
```

**Ready to query?** See [Query Examples](#-query-examples) below.

## 📡 API Overview

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/query` | POST | Execute SQL queries |
| `/health` | GET | Health check |
| `/metrics` | GET | System metrics |

**Complete API documentation:** [API.md](API.md)

## ⚙️ Configuration

Configure GoDuck using environment variables:

| Variable | Default | Description | Valid Range |
|----------|---------|-------------|-------------|
| `GODUCK_DATABASE_PATH` | *Optional* | Path to DuckDB file (uses in-memory if not specified) | Any valid file path or empty |
| `GODUCK_PORT` | `8080` | HTTP server port | 1-65535 |
| `GODUCK_QUERY_TIMEOUT` | `30s` | Query execution timeout | 1s-10m |
| `GODUCK_MAX_CONNECTIONS` | `10` | Database connection pool size | 1-100 |
| `GODUCK_LOG_LEVEL` | `info` | Log level | debug, info, warn, error |
| `GODUCK_READ_WRITE` | `false` | Enable read-write access (required for in-memory databases) | true, false |

### 📋 Common Configurations

**Basic File Database (Recommended for Production):**
```bash
export GODUCK_DATABASE_PATH="/path/to/your/database.duckdb"
export GODUCK_PORT="8080"
./goduck
```

**In-Memory Database (Testing/Development):**
```bash
export GODUCK_READ_WRITE="true"  # Required for in-memory
export GODUCK_PORT="8080"
./goduck
```

**High-Traffic Production:**
```bash
export GODUCK_DATABASE_PATH="/var/lib/goduck/production.duckdb"
export GODUCK_MAX_CONNECTIONS="25"
export GODUCK_QUERY_TIMEOUT="60s"
export GODUCK_LOG_LEVEL="warn"
./goduck
```

## 🐳 Deployment

### Docker (Recommended)
```bash
# Run with your database file
docker run -p 8080:8080 \
  -v /path/to/your/database.duckdb:/data/database.duckdb:ro \
  -e GODUCK_DATABASE_PATH=/data/database.duckdb \
  goduck
```

### Docker Compose
```yaml
version: '3.8'
services:
  goduck:
    image: goduck:latest
    ports:
      - "8080:8080"
    volumes:
      - ./data/sample.duckdb:/data/database.duckdb:ro
    environment:
      - GODUCK_DATABASE_PATH=/data/database.duckdb
      - GODUCK_MAX_CONNECTIONS=20
```

### Production Checklist
- [ ] Set `GODUCK_MAX_CONNECTIONS` based on expected load (default: 10)
- [ ] Configure `GODUCK_QUERY_TIMEOUT` for complex queries (default: 30s)
- [ ] Set `GODUCK_LOG_LEVEL=warn` for production (default: info)
- [ ] Monitor `/metrics` endpoint for performance
- [ ] Set up reverse proxy with HTTPS
- [ ] Monitor `/health` endpoint for availability

## 📖 Query Examples

**See [GETTING_STARTED.md](GETTING_STARTED.md) for more examples and use cases.**

### Basic Queries
```bash
# Count records
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT count(*) as total FROM users"}'

# Filter and sort
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT name, age FROM users WHERE age > 25 ORDER BY age DESC LIMIT 10"}'
```

### Advanced Analytics
```bash
# Aggregation with GROUP BY
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT department, AVG(salary) as avg_salary FROM employees GROUP BY department"}'

# Join operations
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT u.name, COUNT(o.id) as order_count FROM users u LEFT JOIN orders o ON u.id = o.user_id GROUP BY u.id, u.name"}'
```

## 🔒 Security Features

### Read-Only vs Read-Write Database Access
- **File databases**: Opened in **read-only mode** by default unless `GODUCK_READ_WRITE=true` is set
- **In-memory databases**: Require **read-write mode** (`GODUCK_READ_WRITE=true`) when no `GODUCK_DATABASE_PATH` is specified
- File databases in read-only mode block all write operations (INSERT, UPDATE, DELETE, CREATE, DROP)
- Read-write mode allows full SQL operations for testing and development

### Rate Limiting
- **60 requests per minute** per IP address
- Automatic cleanup of inactive connections
- Returns `429 Too Many Requests` when limit exceeded

### Input Validation
- Query size limited to **10KB** to prevent abuse
- SQL injection protection through parameterized queries
- Request timeout enforcement
- Malformed JSON request handling

### Error Handling
- Sanitized error messages (no internal details exposed)
- Request ID tracking for debugging
- Structured logging for security monitoring

## 📈 Monitoring

### Health Check
```bash
curl http://localhost:8080/health
# Returns: {"status": "healthy", "time": "2025-07-21T19:08:27-04:00"}
```

### System Metrics
```bash
curl http://localhost:8080/metrics
# Returns system stats, database pool status, and performance metrics
```

### Key Metrics to Monitor
- **Connection Pool Usage**: Available in `/metrics` - watch for pool exhaustion
- **Query Response Times**: Track via `/metrics` endpoint  
- **Error Rates**: Monitor 4xx/5xx responses
- **Memory Usage**: System memory stats in `/metrics`

## 🚨 Troubleshooting

### Common Issues
| Error | Cause | Solution |
|-------|-------|----------|
| "in-memory database requires read-write access" | No `GODUCK_READ_WRITE=true` set | Set `GODUCK_READ_WRITE=true` for in-memory databases |
| "failed to open database" | Invalid file path or permissions | Check file path and permissions |
| "Query execution failed" | Invalid SQL syntax | Check SQL syntax |
| "Rate limit exceeded" | Too many requests (60/min per IP) | Wait and retry |
| "Query too large" | SQL > 10KB | Reduce query size |

### HTTP Status Codes
- `200` - Success
- `400` - Bad Request (invalid SQL, query too large)
- `429` - Too Many Requests (rate limit exceeded)
- `500` - Internal Server Error (database issues)
- `503` - Service Unavailable (health check failed)

## 🔧 Advanced Configuration

For advanced use cases, see the full configuration options:

| Variable | Default | Description |
|----------|---------|-------------|
| `GODUCK_DATABASE_PATH` | *Optional* | Path to DuckDB file (uses in-memory if not specified) |
| `GODUCK_PORT` | `8080` | HTTP server port |
| `GODUCK_QUERY_TIMEOUT` | `30s` | Query execution timeout |
| `GODUCK_MAX_CONNECTIONS` | `10` | Database connection pool size |
| `GODUCK_LOG_LEVEL` | `info` | Log level (debug, info, warn, error) |
| `GODUCK_READ_WRITE` | `false` | Enable read-write access (required for in-memory) |

### Kubernetes
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: goduck
spec:
  replicas: 3
  selector:
    matchLabels:
      app: goduck
  template:
    spec:
      containers:
      - name: goduck
        image: goduck:latest
        ports:
        - containerPort: 8080
        env:
        - name: GODUCK_DATABASE_PATH
          value: "/data/database.duckdb"
        - name: GODUCK_MAX_CONNECTIONS
          value: "20"
        volumeMounts:
        - name: database
          mountPath: /data
          readOnly: true
      volumes:
      - name: database
        persistentVolumeClaim:
          claimName: database-pvc
```

## 🤝 Support

### Getting Help
- 📖 Check this documentation first
- 🐛 Report bugs via GitHub issues  
- 💡 Request features via GitHub discussions

### Frequently Asked Questions
**Q: Can I modify data through GoDuck?**  
A: File databases are read-only by default. Set `GODUCK_READ_WRITE=true` to enable writes. In-memory databases always require read-write mode.

**Q: What's the maximum query size?**  
A: Queries are limited to 10KB to prevent abuse.

**Q: How do I monitor performance?**  
A: Use the `/metrics` endpoint for system statistics and connection pool status.

---
**Built for the DuckDB community** 🦆