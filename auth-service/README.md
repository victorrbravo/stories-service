# Stories auth service

TypeScript auth service built on [better-auth](https://www.better-auth.com) + Hono. It uses the existing
`stories` Postgres database and shares `ta_user` / `ta_session` with the Go service.

## How it maps onto the existing tables

| better-auth model | Table             | Notes |
|-------------------|-------------------|-------|
| user              | `ta_user`         | `emailVerified` → `activated`, `createdAt` → `created_at`; extra fields `firstname`, `lastname`, `username`, `accountType`, `status` |
| session           | `ta_session`      | `userId` → `fk_user`, `createdAt` → `started_at` |
| account           | `ta_account` (new)      | Credential + social provider accounts |
| verification      | `ta_verification` (new) | Email verification / password reset tokens |

`migrations/001_better_auth.sql` prepares the schema. It is additive and safe to re-run:

- better-auth always addresses rows by a column named `id`, so `ta_user`, `ta_session` (and the new tables) get
  `id GENERATED ALWAYS AS (pk_…) STORED`.
- Adds `ta_user.image`, `ta_user.updated_at`, `ta_session.expires_at/updated_at/ip_address/user_agent`.
- Changes `ta_session.token` from `char(64)` to `varchar(255)` (space-padded `char` values break signed cookies).
- Creates `ta_account` and `ta_verification`, and backfills a `credential` account for every existing user.
- **Password sync triggers:** the Go service reads `ta_user.password`, better-auth reads `ta_account.password`.
  Both use bcrypt (cost 10), and triggers mirror the hash in both directions. A user can sign up or change
  their password in either service and still log in through both.

On sign-up the service sends the Resend template `signup` (variables `first_name`, `code`, `expiration`), the
same email the Go signup sends. `code` is an 8-character token stored in `ta_session` (no `expires_at`, so it is
never a valid session); the echoServer `/confirm` page checks it, emails a 6-digit PIN and sets
`ta_user.activated`, which better-auth reports as `emailVerified`. Needs `API_RESEND_KEY` in `.env`.

On sign-up, `firstname`/`lastname` are derived from `name` and `username` defaults to the email when the
client doesn't send them.

## Run

```bash
cp .env.example .env            # set BETTER_AUTH_SECRET (openssl rand -base64 32)
npm install
npm run build && npm run migrate   # once, applies migrations/*.sql
npm run dev                        # http://localhost:3000
```

With Docker (connects to the Postgres on the host through `host.docker.internal`):

```bash
docker compose up -d --build
docker compose run --rm auth node dist/migrate.js   # once
```

## Endpoints

All better-auth routes are mounted under `/auth` (in production nginx routes `https://doitall.work/auth/` here), plus `GET /health`. Common ones:

| Method | Path | Body |
|--------|------|------|
| POST | `/auth/sign-up/email` | `{ name, email, password, firstname?, lastname?, username? }` |
| POST | `/auth/sign-in/email` | `{ email, password }` |
| GET  | `/auth/get-session` | — |
| POST | `/auth/sign-out` | — |
| POST | `/auth/change-password` | `{ currentPassword, newPassword }` |
| POST | `/auth/request-password-reset` | `{ email, redirectTo }` |
| GET  | `/auth/verify-email?token=…` | — |

Browsers authenticate with the `better-auth.session_token` cookie. Other clients can use the `token` returned
by sign-in (also in the `set-auth-token` header) as `Authorization: Bearer <token>`.

## Environment

See `.env.example`. `BETTER_AUTH_SECRET` is required in production. Without `API_RESEND_KEY`, verification and
reset emails are printed to the console. Set `TRUST_PROXY=true` only behind a reverse proxy that sets
`X-Forwarded-For`.
