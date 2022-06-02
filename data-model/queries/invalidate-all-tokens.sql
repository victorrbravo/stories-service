UPDATE "sessions"
SET "invalid" = TRUE
WHERE "user_id" = $1