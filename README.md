# stories-service

Service for managing user stories and authentication.

## Tech Stack
- **Language**: Go 1.24
- **Package Manager**: Go Modules
- **Key Libraries**: 
  - `lib/pq` (PostgreSQL driver)
  - `jwt-go` (Authentication)
  - `go-openai` (AI integration)
  - `aws-sdk-go` (AWS services)
  - `resend-go` (Email)
  - `d2` (Diagramming)

## Overview
A Go-based service providing authentication and story management capabilities.

## Requirements
- Go 1.24+

## Setup
1. Clone the repository.
2. Install dependencies:
   ```bash
   go mod download
   ```

## Run
To run the authentication service:
```bash
go run cmd/auth/main.go
```

## Scripts
The project uses `go generate` for code generation. Run the following to update generated code:
```bash
go generate ./cmd/auth/...
```

## Environment Variables
The application relies on several environment variables.
- `DATABASE_URL`: Connection string for the PostgreSQL database (TODO: Add more variables)
- `JWT_SECRET`: Secret key for JWT token signing (TODO: Add more variables)

## Testing
To run the tests:
```bash
go test ./...
```

## Project Structure
- `cmd/`: Entry points for services (e.g., `cmd/auth`).
- `handlers/`: HTTP request handlers.
- `pkg/`: Shared packages and utilities.
- `generated/`: Auto-generated code (models, queries, inputs).
- `generators/`: Tools for code generation.
- `database/`: Database connection and abstraction layer.
- `types/`: Custom types and data structures.

## License
TODO: Add license information.
