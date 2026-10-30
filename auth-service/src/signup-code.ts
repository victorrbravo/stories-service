import { randomInt } from "node:crypto";
import { pool } from "./db.js";

// Mirrors the Go signup (handlers/auth/signup.go): an 8-character code stored
// as a ta_session token and emailed through the "signup" template. The link in
// that template goes to the echoServer /confirm page, which checks the code,
// emails a 6-digit PIN and sets ta_user.activated (better-auth's emailVerified).
// The row has no expires_at, so it can never be used as a better-auth session.
const CODE_LENGTH = 8;
const CODE_ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";

// Minutes the code stays valid; echoServer enforces it (confirmEmailCodeTTL).
export const SIGNUP_CODE_TTL_MINUTES = 15;

export async function createSignupCode(userId: number | string): Promise<string> {
  let code = "";
  for (let i = 0; i < CODE_LENGTH; i++) code += CODE_ALPHABET[randomInt(CODE_ALPHABET.length)];
  await pool.query("INSERT INTO ta_session (token, fk_user) VALUES ($1, $2)", [code, Number(userId)]);
  return code;
}
