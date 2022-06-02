SELECT "id",
       "first_name",
       "last_name",
       "email",
       "account_type",
       "password",
       "created_at",
       "activated"
FROM "users"
WHERE "email" = $1
