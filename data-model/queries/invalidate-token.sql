UPDATE "sessions"
SET "invalid" = True
WHERE "user_id" = $1
  AND "token" = $2
