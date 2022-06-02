SELECT "id",
       "first_name",
       "last_name",
       "email",
       "account_type",
       "password",
       "created_at",
       "activated"
FROM public."sessions"
         JOIN public."users" ON public."users"."id" = "user_id"
WHERE public."sessions"."token" = $1
  AND "invalid" = False
