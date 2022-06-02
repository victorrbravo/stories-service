
package queries
// Generated file
// ALL EDITS WILL BE LOST

// Automatically generated queries
const (
	CheckSession = `SELECT "id", "first_name", "last_name", "email", "account_type", "password", "created_at", "activated" FROM public."sessions" JOIN public."users" ON public."users"."id" = "user_id" WHERE public."sessions"."token" = $1 AND "invalid" = False`
	CreateSession = `INSERT INTO "sessions" ("token", "user_id", "started_at")VALUES ($1, $2, NOW())`
	CreateUser = `INSERT INTO "users" ("first_name", "last_name", "password", "email", "account_type", "professional_category", "message_from_user", "pre_registration")VALUES ($1::text, $2::text, $3, $4, $5, $6, $7, $8)RETURNING "id"`
	GetUserByEmail = `SELECT "id", "first_name", "last_name", "email", "account_type", "password", "created_at", "activated" FROM "users" WHERE "email" = $1`
	GetUserByID = `SELECT "id", "first_name", "last_name", "email", "account_type", "password", "created_at", "activated" FROM "users" WHERE "id" = $1`
	GetUserByToken = `SELECT "id", "first_name", "last_name", "email", "account_type", "password", "created_at", "activated" FROM "users" WHERE "token" = $1`
	GetUserSessionCount = `SELECT COUNT(*)FROM "sessions" WHERE "user_id" = $1`
	GetUsers = `SELECT "id", "first_name", "last_name", "email", "account_type", "password", "created_at", "activated" FROM "users"`
	InvalidateAllTokens = `UPDATE "sessions" SET "invalid" = TRUE WHERE "user_id" = $1`
	InvalidateToken = `UPDATE "sessions" SET "invalid" = True WHERE "user_id" = $1 AND "token" = $2`
	UpdatePassword = `UPDATE "users" SET "password" = $1 WHERE "id" = $2`
	UpdateUser = `UPDATE "users" SET "first_name" = COALESCE($1, "first_name"), "last_name" = COALESCE($2, "last_name"), "password" = COALESCE($3, "password"), "account_type" = COALESCE($4, "account_type"), "pre_registration" = COALESCE($5, "pre_registration"), "activated" = COALESCE($6, "activated")WHERE "email" = $7` 
)
