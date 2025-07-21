# Getting Started with GoDuck

Welcome to GoDuck! This guide will help you get up and running quickly.

## 🎯 What is GoDuck?

GoDuck is a REST API server that lets you query DuckDB databases over HTTP. It's perfect for:
- 📊 **Data Analytics**: Query large datasets without complex setups
- 🔍 **Data Exploration**: Interactive SQL queries via HTTP
- 📈 **Dashboards**: Backend for data visualization tools
- 🤖 **Applications**: Integrate analytical queries into your apps

## 🚀 Quick Start (5 minutes)

### Step 1: Get GoDuck
```bash
# Option A: Download binary (recommended)
wget https://github.com/your-org/goduck/releases/latest/download/goduck
chmod +x goduck

# Option B: Build from source
git clone https://github.com/your-org/goduck.git
cd goduck
go build -o goduck
```

### Step 2: Prepare Your Database
```bash
# If you have a DuckDB file already, skip to Step 3
# Otherwise, create a sample database:
./start_example.sh
```

### Step 3: Start the Server
```bash
# Set the path to your DuckDB file
export DATABASE_PATH="/path/to/your/database.duckdb"

# Start GoDuck
./goduck
```

You should see:
```
{"level":"info","msg":"Database connection established","time":"2025-07-21T19:12:14-04:00"}
{"level":"info","msg":"Starting server","port":"8080","time":"2025-07-21T19:12:14-04:00"}
```

### Step 4: Make Your First Query
```bash
# Test the connection
curl http://localhost:8080/health

# Run a simple query
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT 1 as hello, '\''world'\'' as message"}'
```

🎉 **Congratulations!** You're now running GoDuck!

---

## 🏗️ Common Use Cases

### 1. Data Exploration
Perfect for exploring large datasets:
```bash
# See what tables exist
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SHOW TABLES"}'

# Explore table structure
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "DESCRIBE users"}'

# Quick data preview
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT * FROM users LIMIT 5"}'
```

### 2. Analytics Queries
Run complex analytical queries:
```bash
# Aggregations
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT department, AVG(salary) as avg_salary FROM employees GROUP BY department ORDER BY avg_salary DESC"}'

# Time series analysis
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT DATE_TRUNC('\''month'\'', order_date) as month, SUM(amount) as total_sales FROM orders GROUP BY month ORDER BY month"}'
```

### 3. Dashboard Backend
Use as a backend for dashboards:
```javascript
// JavaScript example
async function fetchData(sql) {
  const response = await fetch('http://localhost:8080/query', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ sql })
  });
  
  const data = await response.json();
  return data;
}

// Usage
const users = await fetchData('SELECT count(*) as total_users FROM users');
console.log(`Total users: ${users.rows[0][0]}`);
```

---

## ⚙️ Configuration Guide

**For complete configuration details, see [README.md](README.md#configuration)**

### Environment Variables
```bash
# Required
export DATABASE_PATH="/path/to/your/database.duckdb"

# Optional (with defaults)
export PORT="8080"                    # Server port
export QUERY_TIMEOUT="30s"           # Query timeout
export MAX_CONNECTIONS="10"          # Database connections
export LOG_LEVEL="info"              # Logging level
```

### Performance Tuning

**For Small Workloads (< 10 concurrent users):**
```bash
export MAX_CONNECTIONS="5"
export QUERY_TIMEOUT="30s"
```

**For Medium Workloads (10-50 concurrent users):**
```bash
export MAX_CONNECTIONS="15"
export QUERY_TIMEOUT="60s"
```

**For Large Workloads (50+ concurrent users):**
```bash
export MAX_CONNECTIONS="25"
export QUERY_TIMEOUT="120s"
```

## 🐳 Docker Quick Start

**For complete Docker deployment options, see [README.md](README.md#docker-deployment)**

```bash
# Quick run with your database
docker run -p 8080:8080 \
  -v /path/to/your/database.duckdb:/data/database.duckdb:ro \
  -e DATABASE_PATH=/data/database.duckdb \
  goduck:latest
```

---

## 🔍 Monitoring & Debugging

**For complete monitoring details, see [README.md](README.md#performance--monitoring)**

### Health Checks
```bash
# Basic health check
curl http://localhost:8080/health

# Detailed metrics
curl http://localhost:8080/metrics | jq .
```

### Logs
GoDuck outputs structured JSON logs:
```bash
# View logs with jq for better formatting
./goduck 2>&1 | jq .

# Filter for errors only
./goduck 2>&1 | jq 'select(.level == "error")'
```

### Request Tracing
Add request IDs for debugging:
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: my-debug-request-123" \
  -d '{"sql": "SELECT 1"}'
```

---

## 🚨 Troubleshooting

**For complete troubleshooting guide, see [README.md](README.md#error-handling)**

### Common Issues

**"DATABASE_PATH is required"**
```bash
# Solution: Set the environment variable
export DATABASE_PATH="/path/to/your/database.duckdb"
```

**"failed to open database"**
```bash
# Check if file exists and has correct permissions
ls -la /path/to/your/database.duckdb
chmod 644 /path/to/your/database.duckdb
```

**"Rate limit exceeded"**
```bash
# Wait a minute or reduce request frequency
# Current limit: 60 requests per minute per IP
```

**"Query execution failed"**
```bash
# Check your SQL syntax
# Remember: only SELECT queries allowed
# Use a DuckDB client to test queries first
```

### Getting Help
- 📖 Read the full [README.md](README.md)
- 📋 Check the [API Reference](API.md)
- 🐛 Report issues on GitHub
- 💬 Ask questions in discussions

---

## 🎓 Next Steps

1. **Read the API Documentation**: See [API.md](API.md) for complete endpoint details
2. **Explore Advanced Features**: Check out the metrics endpoint and request tracing in [README.md](README.md)
3. **Production Deployment**: Review security and performance considerations in [README.md](README.md#deployment)
4. **Build Something Cool**: Use GoDuck as a backend for your data applications!

Happy querying! 🚀
