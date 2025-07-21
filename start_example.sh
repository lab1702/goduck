#!/bin/bash

# Example usage script for GoDuck DuckDB REST API Server

echo "🦆 GoDuck - DuckDB REST API Server Example"
echo "==========================================="
echo ""

# Create a sample DuckDB database (requires duckdb CLI)
echo "Creating sample database..."
if command -v duckdb &> /dev/null; then
    duckdb data/sample.duckdb << EOF
CREATE TABLE users AS SELECT 
    id, 
    'user_' || id as name, 
    random() * 100 as score 
FROM range(1, 1001) t(id);

CREATE TABLE orders AS SELECT 
    id,
    (random() * 1000)::INT + 1 as user_id,
    random() * 1000 as amount
FROM range(1, 5001) t(id);
EOF
    echo "Sample database created with users and orders tables"
else
    echo "DuckDB CLI not found. Please install DuckDB to create sample data:"
    echo "  wget https://github.com/duckdb/duckdb/releases/latest/download/duckdb_cli-linux-amd64.zip"
    echo "  unzip duckdb_cli-linux-amd64.zip"
    echo "  sudo mv duckdb /usr/local/bin/"
fi

echo ""
echo "Starting GoDuck server..."
echo "Database: data/sample.duckdb"
echo "Port: 8080"
echo ""

# Set environment variables and start server
export DATABASE_PATH="$(pwd)/data/sample.duckdb"
export PORT="8080"
export LOG_LEVEL="info"
export QUERY_TIMEOUT="30s"
export MAX_CONNECTIONS="10"

echo "Environment:"
echo "  DATABASE_PATH=$DATABASE_PATH"
echo "  PORT=$PORT"
echo "  LOG_LEVEL=$LOG_LEVEL"
echo "  QUERY_TIMEOUT=$QUERY_TIMEOUT"
echo "  MAX_CONNECTIONS=$MAX_CONNECTIONS"
echo ""

echo "🔍 Example queries to test:"
echo ""
echo "1. 💚 Health check:"
echo "   curl http://localhost:8080/health"
echo ""
echo "2. 📊 System metrics:"
echo "   curl http://localhost:8080/metrics"
echo ""
echo "3. 👥 Count users:"
echo '   curl -X POST http://localhost:8080/query \'
echo '     -H "Content-Type: application/json" \'
echo '     -d '"'"'{"sql": "SELECT count(*) as user_count FROM users"}'"'"
echo ""
echo "4. 🏆 Top users by score:"
echo '   curl -X POST http://localhost:8080/query \'
echo '     -H "Content-Type: application/json" \'
echo '     -d '"'"'{"sql": "SELECT name, score FROM users ORDER BY score DESC LIMIT 5"}'"'"
echo ""
echo "5. 📈 User order summary:"
echo '   curl -X POST http://localhost:8080/query \'
echo '     -H "Content-Type: application/json" \'
echo '     -d '"'"'{"sql": "SELECT u.name, COUNT(o.id) as order_count, SUM(o.amount) as total_spent FROM users u LEFT JOIN orders o ON u.id = o.user_id GROUP BY u.id, u.name ORDER BY total_spent DESC LIMIT 10"}'"'"
echo ""
echo "🚀 Starting server (Ctrl+C to stop)..."
echo "📋 Features: Rate limiting (60/min), Request tracing, Query size limits (10KB)"

# Start the server
./goduck