# GoDuck - DuckDB REST API Server

A production-ready, high-performance REST API server for executing SQL queries against read-only DuckDB databases.

## 📚 Documentation

- 🚀 **New to GoDuck?** Start with [GETTING_STARTED.md](GETTING_STARTED.md)
- 📋 **API Reference** [API.md](API.md)  
- 📝 **Version History** [CHANGELOG.md](CHANGELOG.md)
- 📚 **Documentation Guide** [DOCS.md](DOCS.md)

## ✨ Features

- **🔒 Security First**: Read-only database mode with rate limiting and input validation
- **⚡ High Performance**: Connection pooling, query timeouts, and optimized execution
- **📊 Production Ready**: Comprehensive monitoring, metrics, and structured logging
- **🛡️ Robust**: Request tracing, graceful shutdown, and error handling
- **🐳 Container Ready**: Docker support with multi-stage builds
- **📈 Observable**: Health checks, metrics endpoint, and request tracing

## 🚀 Quick Start

**New users:** See [GETTING_STARTED.md](GETTING_STARTED.md) for a detailed walkthrough.

```bash
# 1. Set database path (required)
export DATABASE_PATH="/path/to/your/database.duckdb"

# 2. Run the server
./goduck

# 3. Test it works
curl http://localhost:8080/health
```

**Need sample data?** Run `./start_example.sh` to create a test database.

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
| `DATABASE_PATH` | **Required** | Path to DuckDB file (must be specified) | Any valid file path |
| `PORT` | `8080` | HTTP server port | 1-65535 |
| `QUERY_TIMEOUT` | `30s` | Query execution timeout | 1s-10m |
| `MAX_CONNECTIONS` | `10` | Database connection pool size | 1-100 |
| `LOG_LEVEL` | `info` | Log level | debug, info, warn, error |

### 📋 Configuration Examples

**Development:**
```bash
export DATABASE_PATH="./data/development.duckdb"
export LOG_LEVEL="debug"
export MAX_CONNECTIONS="5"
```

**Production:**
```bash
export DATABASE_PATH="/var/lib/goduck/production.duckdb"
export PORT="8080"
export LOG_LEVEL="info"
export QUERY_TIMEOUT="60s"
export MAX_CONNECTIONS="25"
```

## 🏃 Usage Examples

### Running Locally
```bash
# Build and run
go build -o goduck
DATABASE_PATH=/path/to/your/database.duckdb ./goduck

# Or with go run
DATABASE_PATH=/path/to/your/database.duckdb go run main.go
```

### Docker Deployment
```bash
# Build image
docker build -t goduck .

# Run container with volume mount
docker run -p 8080:8080 \
  -v /path/to/your/database.duckdb:/data/database.duckdb:ro \
  -e DATABASE_PATH=/data/database.duckdb \
  -e LOG_LEVEL=info \
  goduck
```

### Docker Compose
```yaml
version: '3.8'
services:
  goduck:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./data/sample.duckdb:/data/database.duckdb:ro
    environment:
      - DATABASE_PATH=/data/database.duckdb
      - LOG_LEVEL=info
      - MAX_CONNECTIONS=20
```

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

### Read-Only Database Access
- Database opened in **read-only mode** using DuckDB's `access_mode=read_only`
- All write operations (INSERT, UPDATE, DELETE, CREATE, DROP) are blocked by DuckDB
- No risk of data modification through the API

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

## 📈 Performance & Monitoring

### Connection Pooling
- Configurable pool size (default: 10 connections)
- Automatic connection lifecycle management
- Idle connection cleanup
- Pool statistics available via `/metrics`

### Query Performance
- Configurable query timeouts (default: 30s)
- Context-based cancellation
- Execution time tracking
- Connection reuse optimization

### Observability
- **Request Tracing**: Unique request IDs for debugging
- **Structured Logging**: JSON format with contextual information
- **Metrics Endpoint**: System and database statistics
- **Health Checks**: Database connectivity monitoring

### Monitoring Metrics
The `/metrics` endpoint provides:
- **System**: Memory usage, goroutines, GC stats
- **Database**: Connection pool status, query statistics
- **Performance**: Uptime, response times
- **Timestamp**: For time-series monitoring

## 🚨 Error Handling

### HTTP Status Codes
- `200` - Success
- `400` - Bad Request (invalid SQL, query too large)
- `429` - Too Many Requests (rate limit exceeded)
- `500` - Internal Server Error (database issues)
- `503` - Service Unavailable (health check failed)

### Error Response Format
```json
{
  "error": "Human-readable error message",
  "timestamp": "2025-07-21T19:08:27-04:00"
}
```

### Common Issues
| Error | Cause | Solution |
|-------|-------|----------|
| "DATABASE_PATH is required" | Missing env variable | Set `DATABASE_PATH` |
| "Query execution failed" | Invalid SQL | Check SQL syntax |
| "Rate limit exceeded" | Too many requests | Wait and retry |
| "Query too large" | SQL > 10KB | Reduce query size |

## 🛠️ Development

### Project Structure
```
goduck/
├── main.go                 # Application entry point
├── internal/
│   ├── config/            # Configuration management
│   ├── database/          # Database connection handling
│   ├── handlers/          # HTTP request handlers
│   └── middleware/        # HTTP middleware
├── pkg/models/            # Data models
├── tests/                 # Test files
├── data/                  # Sample databases
└── Dockerfile            # Container definition
```

### Building from Source
```bash
# Clone repository
git clone <repository-url>
cd goduck

# Install dependencies
go mod download

# Build binary
go build -o goduck

# Run tests
go test ./...
```

### Contributing
1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## 📦 Deployment

### Production Checklist
- [ ] Set appropriate `MAX_CONNECTIONS` for your load
- [ ] Configure `QUERY_TIMEOUT` based on query complexity
- [ ] Set `LOG_LEVEL=info` or `warn` for production
- [ ] Monitor `/metrics` endpoint
- [ ] Set up log aggregation
- [ ] Configure reverse proxy (nginx, etc.)
- [ ] Enable HTTPS termination
- [ ] Set up monitoring alerts

### Kubernetes Deployment
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
    metadata:
      labels:
        app: goduck
    spec:
      containers:
      - name: goduck
        image: goduck:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_PATH
          value: "/data/database.duckdb"
        - name: MAX_CONNECTIONS
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
- 📧 Contact maintainers for security issues

### Common Questions
**Q: Can I modify the database through GoDuck?**  
A: No, GoDuck enforces read-only access. All write operations are blocked.

**Q: What's the maximum query size?**  
A: Queries are limited to 10KB to prevent abuse.

**Q: How do I monitor performance?**  
A: Use the `/metrics` endpoint and monitor connection pool statistics.

**Q: Can I use with other databases?**  
A: No, GoDuck is specifically designed for DuckDB files.

---

**Made with ❤️ for the DuckDB community**