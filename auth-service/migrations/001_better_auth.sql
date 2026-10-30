-- Adapts the existing ta_user / ta_session tables for better-auth and adds the
-- account / verification tables it needs. Additive and idempotent: existing
-- columns used by the Go service are left untouched (except ta_session.token,
-- see below), so both services can share the same rows.

BEGIN;

-- ---------------------------------------------------------------------------
-- ta_user  (better-auth "user" model)
-- better-auth always addresses rows by a column literally named "id", so we
-- expose pk_user under that name as a generated column.
-- ---------------------------------------------------------------------------
ALTER TABLE ta_user ADD COLUMN IF NOT EXISTS id integer GENERATED ALWAYS AS (pk_user) STORED;
CREATE UNIQUE INDEX IF NOT EXISTS ta_user_id_key ON ta_user (id);
ALTER TABLE ta_user ADD COLUMN IF NOT EXISTS image varchar;
ALTER TABLE ta_user ADD COLUMN IF NOT EXISTS updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP;

-- ---------------------------------------------------------------------------
-- ta_session  (better-auth "session" model)
-- token was char(64): fixed-width chars come back space-padded, which corrupts
-- the signed session cookie. varchar keeps the same values (padding trimmed).
-- ---------------------------------------------------------------------------
ALTER TABLE ta_session ALTER COLUMN token TYPE varchar(255);
ALTER TABLE ta_session ADD COLUMN IF NOT EXISTS id integer GENERATED ALWAYS AS (pk_session) STORED;
CREATE UNIQUE INDEX IF NOT EXISTS ta_session_id_key ON ta_session (id);
CREATE INDEX IF NOT EXISTS ta_session_token_idx ON ta_session (token);
CREATE INDEX IF NOT EXISTS ta_session_fk_user_idx ON ta_session (fk_user);
ALTER TABLE ta_session ADD COLUMN IF NOT EXISTS expires_at timestamptz;
ALTER TABLE ta_session ADD COLUMN IF NOT EXISTS updated_at timestamptz DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE ta_session ADD COLUMN IF NOT EXISTS ip_address varchar;
ALTER TABLE ta_session ADD COLUMN IF NOT EXISTS user_agent varchar;

-- ---------------------------------------------------------------------------
-- ta_account  (better-auth "account" model: credentials + social providers)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ta_account (
    pk_account               serial PRIMARY KEY,
    id                       integer GENERATED ALWAYS AS (pk_account) STORED UNIQUE,
    fk_user                  integer NOT NULL REFERENCES ta_user (pk_user) ON DELETE CASCADE,
    account_id               varchar NOT NULL,
    provider_id              varchar NOT NULL,
    access_token             text,
    refresh_token            text,
    id_token                 text,
    access_token_expires_at  timestamptz,
    refresh_token_expires_at timestamptz,
    scope                    text,
    password                 text,
    created_at               timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at               timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT ta_account_user_provider_key UNIQUE (fk_user, provider_id)
);

-- ---------------------------------------------------------------------------
-- ta_verification  (better-auth "verification" model: email / reset tokens)
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS ta_verification (
    pk_verification serial PRIMARY KEY,
    id              integer GENERATED ALWAYS AS (pk_verification) STORED UNIQUE,
    identifier      varchar NOT NULL,
    value           text NOT NULL,
    expires_at      timestamptz NOT NULL,
    created_at      timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS ta_verification_identifier_idx ON ta_verification (identifier);

-- ---------------------------------------------------------------------------
-- Password sync. The Go service reads/writes ta_user.password, better-auth
-- reads/writes ta_account.password (provider 'credential'). Both use bcrypt,
-- so mirroring the hash keeps login working in either service. The
-- IS DISTINCT FROM guards stop the two triggers from recursing.
-- ---------------------------------------------------------------------------
CREATE OR REPLACE FUNCTION fn_sync_user_password_to_account() RETURNS trigger AS $$
BEGIN
    IF NEW.password IS NULL OR NEW.password = '' THEN
        RETURN NEW;
    END IF;
    INSERT INTO ta_account (fk_user, account_id, provider_id, password)
    VALUES (NEW.pk_user, NEW.pk_user::text, 'credential', NEW.password)
    ON CONFLICT (fk_user, provider_id) DO UPDATE
        SET password = EXCLUDED.password, updated_at = CURRENT_TIMESTAMP
        WHERE ta_account.password IS DISTINCT FROM EXCLUDED.password;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION fn_sync_account_password_to_user() RETURNS trigger AS $$
BEGIN
    IF NEW.provider_id = 'credential' AND NEW.password IS NOT NULL THEN
        UPDATE ta_user SET password = NEW.password
        WHERE pk_user = NEW.fk_user AND password IS DISTINCT FROM NEW.password;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tg_ta_user_password_sync ON ta_user;
CREATE TRIGGER tg_ta_user_password_sync
    AFTER INSERT OR UPDATE OF password ON ta_user
    FOR EACH ROW EXECUTE FUNCTION fn_sync_user_password_to_account();

DROP TRIGGER IF EXISTS tg_ta_account_password_sync ON ta_account;
CREATE TRIGGER tg_ta_account_password_sync
    AFTER INSERT OR UPDATE OF password ON ta_account
    FOR EACH ROW EXECUTE FUNCTION fn_sync_account_password_to_user();

-- Backfill credential accounts for users created by the Go service.
INSERT INTO ta_account (fk_user, account_id, provider_id, password)
SELECT pk_user, pk_user::text, 'credential', password
FROM ta_user
WHERE password IS NOT NULL AND password <> ''
ON CONFLICT (fk_user, provider_id) DO NOTHING;

COMMIT;
