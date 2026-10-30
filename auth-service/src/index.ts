import { serve } from "@hono/node-server";
import { getConnInfo } from "@hono/node-server/conninfo";
import { Hono } from "hono";
import { cors } from "hono/cors";
import { logger } from "hono/logger";
import { auth } from "./auth.js";
import { config } from "./config.js";
import { pool } from "./db.js";

const app = new Hono();

app.use(logger());
app.use(
  "/auth/*",
  cors({
    origin: config.trustedOrigins,
    allowHeaders: ["Content-Type", "Authorization"],
    allowMethods: ["GET", "POST", "OPTIONS"],
    exposeHeaders: ["set-auth-token"],
    credentials: true,
  }),
);

app.get("/health", async (c) => {
  await pool.query("SELECT 1");
  return c.json({ status: "ok" });
});

app.on(["GET", "POST"], "/auth/*", (c) => {
  // better-auth rate-limits and records sessions by client IP, read from
  // X-Forwarded-For. Only trust a client-supplied value behind a proxy.
  if (config.trustProxy) return auth.handler(c.req.raw);
  const headers = new Headers(c.req.raw.headers);
  headers.set("x-forwarded-for", getConnInfo(c).remote.address ?? "");
  return auth.handler(new Request(c.req.raw, { headers }));
});

const server = serve({ fetch: app.fetch, port: config.port }, (info) => {
  console.log(`Auth service listening on http://localhost:${info.port}`);
});

const shutdown = () => {
  server.close(() => {
    pool.end().finally(() => process.exit(0));
  });
};
process.on("SIGINT", shutdown);
process.on("SIGTERM", shutdown);
