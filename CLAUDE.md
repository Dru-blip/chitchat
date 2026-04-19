# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Chitchat is a full-stack chat application with magic-link authentication and end-to-end encryption via prekeys.

- **Backend**: Go 1.25 with Echo v5, PostgreSQL (pgx/v5), sqlc for type-safe SQL
- **Frontend**: Vue 3 + TypeScript in `gossip/` directory, Vite, Tailwind CSS v4, Reka UI

## Development Commands

```bash
# Backend - requires DATABASE_URL env var
go run cmd/main.go              # Starts API server on :5050

# Frontend (from gossip/ directory)
cd gossip && npm run dev        # Vite dev server on :5173
cd gossip && npm run build      # Production build
cd gossip && npm run type-check # TypeScript checking

# Database
make sqlc                       # Regenerate sqlc types after schema changes
```

## Architecture

### Backend Structure

**Layer Pattern**: Handler → Service → Repository (sqlc)

- `cmd/api/api.go` - Server initialization, middleware stack, route registration
- `cmd/main.go` - Entry point, loads .env, connects to DB, starts server
- `internal/db/database.go` - Store struct wraps pgxpool + sqlc Queries
- `internal/db/sqlc/` - **Generated code** from sqlc - never edit manually
- `internal/db/migrations/` - SQL migration files
- `internal/db/queries/` - SQL query files for sqlc generation

**Domain Packages**:
- `auth/` - Magic link authentication flow, session middleware
- `users/` - User management
- `keys/` - Prekey management for E2E encryption
- `mailer/` - Async email sending
- `utils/` - Validator, error handler, SHA256 helper

### Service Interface Pattern

Each domain exposes a Service interface with constructor `NewXxx()`:

```go
type Service interface {
    Method(ctx context.Context, ...) (...)
}

func NewService(repo Repository, mailer Mailer) Service {
    return &service{repo: repo, mailer: mailer}
}
```

### Database & sqlc Configuration

See `sqlc.yaml` for type overrides:
- `uuid` → `github.com/google/uuid.UUID`
- `timestamptz` → `*time.Time` (pointer for nullable fields)
- `timestamp` → `time.Time` or `*time.Time`

Store pattern in `internal/db/database.go`:
```go
type Store struct {
    Db      *pgxpool.Pool
    Queries *sqlc.Queries
}
```

### Authentication Flow

1. User submits email → `SendMagicLink()` creates session with SHA256-hashed token
2. Magic link emailed with raw token (URL format: `/verify-link?id={uuid}&token={raw}`)
3. Client verifies → `VerifyMagicLink()` validates token, marks session used
4. Session cookie created via scs (cookie name: `"chisession"`)

Session middleware (`auth.NewSessionMiddleware`) loads user from session into context.

### Error Handling

- **Global handler**: `internal/utils/errors.go` - unwraps Echo responses, converts validator errors to 422 with field details
- **Auth errors**: Defined in `internal/auth/errors.go` as sentinel errors (`ErrInvalidMagicLink`, `ErrMagicLinkExpired`)
- **HTTP errors**: Use `echo.NewHTTPError()` in handlers, global handler converts to JSON

Validator errors produce response:
```json
{"code": 422, "message": "Invalid Input", "details": {"email": "email is required"}}
```

### Email

Async fire-and-forget pattern in `mailer.go`. Requires env vars: `SMTP_HOST`, `SMTP_USER`, `SMTP_PASS`, `SMTP_FROM`, `SMTP_PORT`.

### Frontend (gossip/)

- **Tailwind v4**: Uses `@tailwindcss/vite` plugin
- **UI components**: Reka UI primitives in `src/components/ui/`
- **Class merging**: Use `cn()` from `src/lib/utils.ts`
- **Build tool**: Vite with Vue plugin

## Environment Variables

```bash
DATABASE_URL=postgresql://...
SMTP_HOST=smtp.example.com
SMTP_PORT=465
SMTP_USER=...
SMTP_PASS=...
SMTP_FROM=noreply@example.com
SMTP_SECURE=true
```

## Code Style

### Go
- Interface naming: `Service`, `Repository`, `Mailer` (not `IService`)
- Constructor pattern: `NewXxx()` returns interface, takes dependencies
- Context first: All service methods accept `context.Context`
- Validator tags: `validate:"required,email"` on struct fields

### Vue/TypeScript
- Components use `<script setup lang="ts">`
- Props/Emits: Use TypeScript interfaces
- UI components export via `index.ts` barrel files
