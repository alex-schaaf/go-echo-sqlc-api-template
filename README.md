# go-echo-sqlc-api-template

This template provides a robust starting point fœor building RESTful APIs using the Echo framework and SQLC for database interactions in Go. It includes a well-structured project layout, essential middleware, and example handlers for authentication and user management.

##  Key Features
- [Echo Framework](https://echo.labstack.com/): Fast and scalable web framework for building APIs.
- [SQLC](https://sqlc.dev/): Generate type-safe Go code from SQL queries.
- [SQLite](https://www.sqlite.org/): 
- Database Migrations: Easily manage database schema changes with migration scripts.
- Authentication: JWT-based authentication with login, logout, and registration endpoints (via HTTP-only cookie).
- Configuration Management: Centralized configuration handling.
- Testing Utilities: Setup for in-memory database testing.

## Project Structure
- `api/`: Contains handlers for authentication and user management.
- `cmd/api/`: Entry point for the API server.
- `db/`: Database-related code
    - `migrations/`: Contains the up/down database migration files.
    - `sql`: Contains the `schema.sql` and `query.sql` used to generate SQLC code.
- `lib/`: Libraries for authentication, configuration, and middleware.
    - `test/`: Utilities for setting up tests with an in-memory SQLite database.
