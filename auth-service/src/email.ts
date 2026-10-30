import { config } from "./config.js";

type Email = { to: string; subject: string; html: string };
type TemplateEmail = { to: string; template: string; variables: Record<string, string> };

// Uses the same Resend account as the Go service. Without API_RESEND_KEY the
// email is logged instead, which is enough for local development.
export async function sendEmail({ to, subject, html }: Email): Promise<void> {
  if (!config.resendApiKey) {
    console.log(`[email] to=${to} subject="${subject}"\n${html}`);
    return;
  }
  await postToResend({ from: config.emailFrom, to: [to], subject, html });
}

// Sends one of the templates stored in Resend, referenced by id or alias.
export async function sendTemplateEmail({ to, template, variables }: TemplateEmail): Promise<void> {
  if (!config.resendApiKey) {
    console.log(`[email] to=${to} template=${template} variables=${JSON.stringify(variables)}`);
    return;
  }
  const { id } = await postToResend({
    from: config.emailFrom,
    to: [to],
    template: { id: template, variables },
  });
  console.log(`[email] sent template=${template} to=${to} resend_id=${id}`);
}

async function postToResend(body: object): Promise<{ id: string }> {
  const response = await fetch("https://api.resend.com/emails", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${config.resendApiKey}`,
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });
  if (!response.ok) {
    throw new Error(`Resend error ${response.status}: ${await response.text()}`);
  }
  return (await response.json()) as { id: string };
}
