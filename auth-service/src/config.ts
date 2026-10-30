const env = (name: string, fallback?: string): string => {
  const value = process.env[name] ?? fallback;
  if (value === undefined) throw new Error(`Missing required env var ${name}`);
  return value;
};

const isProduction = process.env.NODE_ENV === "production";

export const config = {
  port: Number(env("PORT", "3000")),
  baseURL: env("BETTER_AUTH_URL", "http://localhost:3000"),
  // Must be set in production; the dev fallback only exists so `npm run dev` works out of the box.
  secret: isProduction
    ? env("BETTER_AUTH_SECRET")
    : env("BETTER_AUTH_SECRET", "dev-only-secret-change-me-dev-only-secret"),
  trustedOrigins: env("TRUSTED_ORIGINS", "http://localhost:5173")
    .split(",")
    .map((origin) => origin.trim())
    .filter(Boolean),
  trustProxy: env("TRUST_PROXY", "false") === "true",
  requireEmailVerification: env("REQUIRE_EMAIL_VERIFICATION", "false") === "true",
  resendApiKey: process.env.API_RESEND_KEY,
  emailFrom: env("EMAIL_FROM", "DoItAll <support@doitall.work>"),
  // Resend template (id or alias) used for the sign-up email.
  signupTemplate: env("RESEND_SIGNUP_TEMPLATE", "signup"),
};
