INSERT INTO "users" ("first_name",
                     "last_name",
                     "password",
                     "email",
                     "account_type",
                     "professional_category",
                     "message_from_user",
                     "pre_registration")
VALUES ($1::text, $2::text, $3, $4, $5, $6, $7, $8)
RETURNING "id"
