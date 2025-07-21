# Changelog

All notable changes to GoDuck will be documented in this file.

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
