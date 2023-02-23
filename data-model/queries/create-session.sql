INSERT INTO "sessions" ("token", "user_id", "started_at")
VALUES ($1, $2, NOW())