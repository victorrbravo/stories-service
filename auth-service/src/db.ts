import pg from "pg";

// Kept out of config.ts so the migrate script doesn't require the auth secret.
export const pool = new pg.Pool({
  connectionString: process.env.DATABASE_URL ?? "postgres://stories:stories@localhost:5432/stories",
});
