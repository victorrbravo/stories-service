UPDATE "users" SET "first_name" = COALESCE($1, "firstname"),
"last_name" = COALESCE($2, "last_name"), 
"password" = COALESCE($3, "password"),
"account_type" = COALESCE($4, "account_type"),
"pre_registration" = COALESCE($5, "pre_registration"),
"activated" = COALESCE($6, "activated")
WHERE "email" = $7