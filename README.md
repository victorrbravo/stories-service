# stories-service

[![Go Reference](https://pkg.go.dev/badge/github.com/victorrbravo/stories-service.svg)](https://pkg.go.dev/github.com/victorrbravo/stories-service)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

A Go library and REST service for user authentication (signup, login, email
verification, password reset, JWT refresh), user management, and a story /
workflow engine backed by PostgreSQL.

You can use it in two ways:

1. **As a library** — import its packages (router, database layer, models,
   workflow engine) into your own Go program.
2. **As a standalone service** — run the bundled HTTP server.

```
go get github.com/victorrbravo/stories-service
```

---

## Requirements

- Go **1.24+**
- PostgreSQL with the **PostGIS** extension enabled
- (Optional) A [Resend](https://resend.com) API key for transactional email

---

## Installation

Add the module to your project:

```bash
go get github.com/victorrbravo/stories-service@latest
```

Or pin a specific release:

```bash
go get github.com/victorrbravo/stories-service@v0.1.0
```

Then import the packages you need:

```go
import (
    "github.com/victorrbravo/stories-service/database"
    "github.com/victorrbravo/stories-service/web"
)
```

---

## Usage as a library

### 1. Mount the full API in your own HTTP server

`web.CreateRouter` returns an `http.Handler` with all routes (`/auth`, `/users`,
`/social`) registered, so you can mount it alongside your own handlers.

```go
package main

import (
    "log"
    "net/http"
    "os"

    "github.com/victorrbravo/stories-service/database"
    "github.com/victorrbravo/stories-service/web"
)

func main() {
    config := database.DatabaseConfig{
        Host:     "localhost",
        Port:     5432,
        Name:     "stories",
        User:     "stories",
        Password: os.Getenv("DB_PASS"),
    }

    var pg database.PostgresQuery
    pg.Initialize(&config)
    if err := pg.Connect(); err != nil {
        log.Fatalf("database connection failed: %v", err)
    }
    query := pg.NewQuery()

    mux := http.NewServeMux()
    mux.Handle("/", web.CreateRouter(query))           // stories-service API
    mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
        w.Write([]byte("ok"))
    })

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### 2. Use the built-in server

```go
server, err := web.NewServer()
if err != nil {
    log.Fatal(err)
}
server.Listen("localhost", 8080, "", query)
```

### 3. Use the database layer directly

All data access goes through the `database.IQuery` interface
(`QueryOne`, `QueryMany`, `QueryManyPaged`, `Exec`, `Commit`). The
PostgreSQL implementation is `database.PostgresQuery`; you can supply your own
implementation (e.g. a mock in tests).

```go
var q database.IQuery = pg.NewQuery()
err := q.Exec(`UPDATE ta_user SET firstname = $1 WHERE pk_user = $2`, "Ada", 42)
```

### 4. Run workflow actions

`input.ProcessInput` executes YAML-described workflow actions. A complete,
runnable program lives in [examples/show-flow-data](examples/show-flow-data):

```go
order := []byte(`
action: show_flow_data
ticket: 243
by_day: true
page_number: 1
page_size: 12
`)
response, err := input.ProcessInput(context.Background(), order, &query)
```

### Main packages

| Package | Purpose |
|---|---|
| `web` | Router factory (`CreateRouter`) and simple server (`NewServer`, `Listen`) |
| `database` | `IQuery` interface, `PostgresQuery` implementation, `DatabaseConfig` |
| `handlers/auth` | Login, signup, verify-email, refresh-token, change-password, resend-code |
| `handlers/social` | `POST /social/add-step/{nameId}` workflow endpoint |
| `handlers_utils` | Custom router and JWT request authorizer |
| `generated/models` | Struct types for all tables (`TaUser`, `TaStory`, …) |
| `generated/inputs` | Request input types with validation |
| `generated/queries` | Generated SQL queries |
| `pkg/workflow` | Workflow execution engine with D2 diagram rendering |
| `types` | Nullable wrappers (`types.String`, `types.Int`, …) and enums |
| `scanners` | Database row scanners for the custom types |

---

## HTTP endpoints

| Method | Path | Auth |
|---|---|---|
| `POST` | `/auth/login` | No |
| `POST` | `/social/add-step/{nameId}` | JWT |
| `GET` | `/users/` | JWT |
| `POST` | `/users/` | JWT |
| `DELETE` | `/users/` | JWT |

Authenticated routes expect an `Authorization: Bearer <access-token>` header.

---

## Configuration

Configuration is read from environment variables. For local development,
copy the template and fill in your own values — `cmd/auth` and the tests load
`.env` automatically, and real environment variables take precedence:

```bash
cp .env.example .env
```

| Variable | Purpose |
|---|---|
| `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASS` | PostgreSQL connection used by `cmd/auth` and the tests |
| `REMOTE_DB_HOST`, `REMOTE_DB_PASS` | Remote database used by some integration tests |
| `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` | AWS credentials (S3 presigned URLs) |
| `JWT_ACCESS_SECRET` | Secret used to sign access tokens |
| `JWT_REFRESH_SECRET` | Secret used to sign refresh tokens |
| `API_RESEND_KEY` | Resend API key for sending email |
| `URL_VERIFY_EMAIL` | Link included in verification emails |
| `URL_CHANGE_PASSWORD` | Link included in password-reset emails |
| `URL_RESEND_EMAIL` | Link for resending the verification code |
| `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD` | SMTP settings (alternative mailer) |
| `OPEN_API_KEY` | OpenAI API key (AI-assisted story features) |
| `GOOGLE_API_KEY` | Google Gemini API key (AI-assisted story features) |
| `WORKFLOW_URL` | Base URL used by the workflow engine |
| `SHOW_EMPTY_TASKS` | Show empty tasks in workflow output |

> If `JWT_ACCESS_SECRET` / `JWT_REFRESH_SECRET` are not set, a random secret is
> generated per process, so tokens are invalidated on restart. Always set them
> in production. Never commit a `.env` file.

---

## Running locally (standalone service)

```bash
git clone https://github.com/victorrbravo/stories-service.git
cd stories-service

# Start PostgreSQL (with PostGIS) and Adminer
docker compose -f docker/docker-compose.yml up -d

# Start the server on localhost:8080
go run cmd/auth/main.go
```

`cmd/auth/main.go` reads the database settings from `DB_*` variables and
defaults to the Docker Compose credentials (`stories` / `stories`).

---

## Development

### Tests

```bash
go test ./...
go test ./handlers/auth/...   # single package
```

Some tests require a running local PostgreSQL instance.

### Code generation

Everything under `generated/` is produced from the JSON specs in
`data-model/` by the generators in `generators/`. **Do not edit generated
files by hand** — edit the specs or templates and re-run:

```bash
go generate ./cmd/auth/...
```

Generated code is committed to the repository so that `go get` works without
running the generators.

### Releasing a new version

Go modules are versioned with git tags following [semver](https://semver.org):

```bash
git tag v0.1.0
git push origin v0.1.0
```

The version then becomes available via `go get` and on
[pkg.go.dev](https://pkg.go.dev/github.com/victorrbravo/stories-service).

---

## Project layout

```
cmd/auth/          entry point of the standalone server
web/               router and HTTP server
handlers/          hand-written handlers (auth, social, users)
handlers_utils/    custom router + JWT middleware
database/          IQuery interface and PostgreSQL implementation
generated/         generated models, inputs, queries and handlers
generators/        code generators (go generate)
data-model/        JSON specs used by the generators
pkg/workflow/      workflow engine + D2 rendering
types/, scanners/  nullable types, enums and DB scanners
docker/            docker-compose for local PostgreSQL
```

---

## Contributing

Issues and pull requests are welcome. Please run `go test ./...` and
`go vet ./...` before submitting.

## License

Copyright 2026 Victor Bravo.

Licensed under the [Apache License, Version 2.0](LICENSE). You may use,
modify and distribute this software, including commercially, provided you
keep the license and the [NOTICE](NOTICE) file. Contributions are accepted
under the same license.
