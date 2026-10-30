import { readdir, readFile } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { pool } from "./db.js";

// Applies every migrations/*.sql file in name order. The files are idempotent,
// so re-running is safe.
const migrationsDir = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../migrations");

const files = (await readdir(migrationsDir)).filter((f) => f.endsWith(".sql")).sort();
try {
  for (const file of files) {
    console.log(`Applying ${file}`);
    await pool.query(await readFile(path.join(migrationsDir, file), "utf8"));
  }
  console.log("Migrations applied");
} finally {
  await pool.end();
}
