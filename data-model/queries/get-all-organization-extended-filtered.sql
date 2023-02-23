WITH CTE3 AS (
    WITH CTE2 AS (
        WITH CTE1 AS (
            WITH CTE AS (
                SELECT "core"."organizations"."id",
                       "core"."organizations"."organization"  AS "name",
                       "ext"."ticker",
                       "ext"."description",
                       "ext"."short_name",
                       "ext"."market_change_percent",
                       "ext"."market_price",
                       "ext"."currency",
                       "ext"."updated_at",
                       COUNT("organization_user_relationships"."portfolio_id")
                       OVER (PARTITION BY "organization_user_relationships"."organization_id") AS "portfolio_count"
                FROM "core"."organizations"
                         LEFT JOIN core."company_data" ext
                                   ON "core"."organizations"."id" = "ext"."fk_organization_id"
                         LEFT JOIN
                     "organization_user_relationships" ON
                             "core"."organizations"."id" = "organization_user_relationships".organization_id
                WHERE "core"."organizations"."organization_type" = 'Publicly Traded Company'
                GROUP BY "core"."organizations"."id",
                         "ext"."description",
                         "ext"."short_name",
                         "ext"."market_change_percent",
                         "ext"."market_price",
                         "ext"."currency",
                         "ext"."updated_at",
                         "ext"."ticker",
                         "organization_user_relationships"."portfolio_id",
                         "organization_user_relationships"."organization_id"
            )
            SELECT "id",
                   "name",
                   "ticker",
                   "description",
                   "short_name",
                   "market_change_percent",
                   "market_price",
                   "currency",
                   "updated_at",
                   "portfolio_count"
            FROM CTE
            WHERE CASE WHEN $1::VARCHAR[] IS NOT NULL THEN "ticker"=ANY($1::VARCHAR[]) ELSE TRUE END
            GROUP BY "id",
                     "name",
                     "ticker",
                     "description",
                     "short_name",
                     "market_change_percent",
                     "market_price",
                     "currency",
                     "updated_at",
                     "portfolio_count"
            ORDER BY "name"
        )
        SELECT "id",
               "name",
               "ticker",
               "description",
               "short_name",
               "market_change_percent",
               "market_price",
               "currency",
               "updated_at",
               "portfolio_count"
        FROM CTE1
        WHERE CASE WHEN $2::VARCHAR[] IS NOT NULL THEN "currency"=ANY($2::VARCHAR[])  ELSE TRUE END
        GROUP BY "id",
                 "name",
                 "ticker",
                 "description",
                 "short_name",
                 "market_change_percent",
                 "market_price",
                 "currency",
                 "updated_at",
                 "portfolio_count"
        ORDER BY "name"
    )
    SELECT "id",
           "name",
           "ticker",
           "description",
           "short_name",
           "market_change_percent",
           "market_price",
           "currency",
           "updated_at",
           "portfolio_count"
    FROM CTE2
    WHERE CASE WHEN $3::VARCHAR[] IS NULL THEN TRUE  ELSE
            (CASE  WHEN "market_change_percent" IS NULL THEN 'Neutral'
                  WHEN "market_change_percent" > 0.0 THEN 'Gain'
                  WHEN "market_change_percent" = 0.0 THEN 'Neutral'
                  WHEN "market_change_percent" < 0.0 THEN 'Loss'
                END) = ANY($3::VARCHAR[])
        END
    GROUP BY "id",
             "name",
             "ticker",
             "description",
             "short_name",
             "market_change_percent",
             "market_price",
             "currency",
             "updated_at",
             "portfolio_count"
    ORDER BY "name"
)
SELECT "id",
       "name",
       "ticker",
       "description",
       "short_name",
       "market_change_percent",
       "market_price",
       CASE WHEN "currency" is null then 'USD' ELSE "currency" END,
       "updated_at",
       "portfolio_count"
FROM CTE3
WHERE CASE
          WHEN $4::VARCHAR[] IS NOT NULL THEN "ticker"=ANY($4::VARCHAR[])
          ELSE  TRUE END
GROUP BY "id",
         "name",
         "ticker",
         "description",
         "short_name",
         "market_change_percent",
         "market_price",
         "currency",
         "updated_at",
         "portfolio_count"
ORDER BY "name"