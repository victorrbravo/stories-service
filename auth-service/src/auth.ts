import { betterAuth } from "better-auth";
import { queueAfterTransactionHook } from "@better-auth/core/context";
import { bearer } from "better-auth/plugins";
import bcrypt from "bcryptjs";
import { config } from "./config.js";
import { pool } from "./db.js";
import { sendEmail, sendTemplateEmail } from "./email.js";
import { createSignupCode, SIGNUP_CODE_TTL_MINUTES } from "./signup-code.js";

// Same algorithm and cost as the Go service (bcrypt.GenerateFromPassword with
// DefaultCost), so hashes are interchangeable between both services.
const BCRYPT_COST = 10;

export const auth = betterAuth({
  appName: "Stories",
  baseURL: config.baseURL,
  // Served at https://doitall.work/auth/* (nginx routes /auth/ to this service).
  basePath: "/auth",
  secret: config.secret,
  trustedOrigins: config.trustedOrigins,
  database: pool,

  advanced: {
    // pk_user / pk_session are serial integers; ids are produced by Postgres.
    database: { generateId: "serial" },
  },

  user: {
    modelName: "ta_user",
    fields: {
      emailVerified: "activated",
      createdAt: "created_at",
      updatedAt: "updated_at",
    },
    additionalFields: {
      firstname: { type: "string", required: false },
      lastname: { type: "string", required: false },
      username: { type: "string", required: false },
      accountType: {
        type: "string",
        fieldName: "account_type",
        required: false,
        input: false,
        defaultValue: "user",
      },
      status: { type: "string", required: false, input: false, defaultValue: "" },
      // NOT NULL column read by the Go service. Filled with the real hash by the
      // ta_account -> ta_user sync trigger once the credential account exists.
      password: {
        type: "string",
        required: false,
        input: false,
        returned: false,
        defaultValue: "",
      },
    },
  },

  session: {
    modelName: "ta_session",
    fields: {
      userId: "fk_user",
      createdAt: "started_at",
      expiresAt: "expires_at",
      updatedAt: "updated_at",
      ipAddress: "ip_address",
      userAgent: "user_agent",
    },
  },

  account: {
    modelName: "ta_account",
    fields: {
      userId: "fk_user",
      accountId: "account_id",
      providerId: "provider_id",
      accessToken: "access_token",
      refreshToken: "refresh_token",
      idToken: "id_token",
      accessTokenExpiresAt: "access_token_expires_at",
      refreshTokenExpiresAt: "refresh_token_expires_at",
      createdAt: "created_at",
      updatedAt: "updated_at",
    },
  },

  verification: {
    modelName: "ta_verification",
    fields: {
      expiresAt: "expires_at",
      createdAt: "created_at",
      updatedAt: "updated_at",
    },
  },

  emailAndPassword: {
    enabled: true,
    requireEmailVerification: config.requireEmailVerification,
    password: {
      hash: (password) => bcrypt.hash(password, BCRYPT_COST),
      verify: ({ hash, password }) => bcrypt.compare(password, hash),
    },
    sendResetPassword: async ({ user, url }) => {
      await sendEmail({
        to: user.email,
        subject: "Reset your password",
        html: `<p>Hi ${user.name},</p><p><a href="${url}">Click here to reset your password</a>.</p>`,
      });
    },
  },

  emailVerification: {
    sendOnSignUp: true,
    // Same email the Go signup sends, so the existing echoServer /confirm flow
    // activates the account. better-auth's own verify link is not used.
    // Sign-up creates the user inside a transaction; the code row references
    // it, so wait for the commit before inserting it.
    sendVerificationEmail: async ({ user }) => {
      await queueAfterTransactionHook(
        async () => {
          const code = await createSignupCode(user.id);
          await sendTemplateEmail({
            to: user.email,
            template: config.signupTemplate,
            variables: {
              first_name: (user as { firstname?: string }).firstname || user.name,
              code,
              expiration: String(SIGNUP_CODE_TTL_MINUTES),
            },
          });
        },
        { onError: (error) => console.error(`[email] sign-up email to ${user.email} failed:`, error) },
      );
    },
  },

  databaseHooks: {
    user: {
      create: {
        // ta_user requires firstname, lastname and a unique username; derive
        // them from the better-auth fields when the client doesn't send them.
        before: async (user) => {
          const [first, ...rest] = (user.name ?? "").trim().split(/\s+/);
          return {
            data: {
              ...user,
              firstname: user.firstname || first || "",
              lastname: user.lastname || rest.join(" "),
              username: user.username || user.email,
            },
          };
        },
      },
    },
  },

  // Lets non-browser clients authenticate with `Authorization: Bearer <token>`.
  plugins: [bearer()],
});
