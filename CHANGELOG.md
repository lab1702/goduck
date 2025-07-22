# Changelog

All notable changes to GoDuck will be documented in this file.

## [0.0.2] - 2025-07-21

### Added
- 🧠 **In-Memory Database Support**: Optional in-memory database mode when DATABASE_PATH is not specified
  - Read-write access for in-memory databases (vs read-only for file databases)
  - Perfect for development, testing, and temporary data analysis
  - Automatic detection: file path provided = read-only file mode, no path = read-write memory mode

### Changed
- 📝 **Configuration**: DATABASE_PATH is now optional (was previously required)
- 📚 **Documentation**: Updated all docs to reflect in-memory database option
- ⬆️ **DuckDB Engine**: Updated from v1.1.3 to v1.3.2 for latest features and performance improvements
- 📦 **Dependencies**: Migrated from go-duckdb v1.8.5 to v2.3.3 for better compatibility and future support

## [0.0.1] - 2025-07-21 - Initial Release

### Added
- �️ **DuckDB Integration**: Read-only access to DuckDB databases
- 🌐 **REST API**: POST /query, GET /health, GET /metrics endpoints
- 🔒 **Security Features**:
  - Read-only database mode enforcement
  - Rate limiting (60 requests per minute per IP)
  - Query size limits (10KB maximum)
  - Input validation and sanitized error messages
- 🏊 **Connection Pooling**: Configurable database connection pool
- ⏱️ **Query Timeouts**: Configurable query execution timeouts
- � **Monitoring & Observability**:
  - Health checks and system metrics
  - Request tracing with unique request IDs
  - Structured JSON logging
  - Connection pool statistics
- � **Container Support**: Multi-stage Dockerfile for production deployment
- 🛑 **Graceful Shutdown**: Proper cleanup on termination
- 🔄 **CORS Support**: Cross-origin request handling
- �️ **Error Handling**: Panic recovery and comprehensive error responses
- ⚙️ **Configuration**: Environment variable based configuration with validation

### Security
- Database opened in read-only mode prevents data modification
- Rate limiting prevents DoS attacks
- Query size limits prevent abuse
- Error messages sanitized to prevent information leakage

### Performance
- Connection pooling for efficient database access
- Context-based query cancellation
- Optimized error handling and response serialization

### Documentation
- 📖 Comprehensive README with production deployment guide
- � Quick start guide (GETTING_STARTED.md)
- � Complete API reference (API.md)
- � Documentation navigation guide (DOCS.md)
