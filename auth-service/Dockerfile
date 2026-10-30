FROM node:22-alpine AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY tsconfig.json ./
COPY src ./src
RUN npm run build && npm prune --omit=dev

FROM node:22-alpine
WORKDIR /app
ENV NODE_ENV=production PORT=3000
# Sign-up email: Resend template "signup" (same one the Go service sends).
# API_RESEND_KEY is a secret, so it is not baked into the image: it comes from
# .env through docker compose, or `docker run -e API_RESEND_KEY=...`.
ENV EMAIL_FROM="DoItAll <support@doitall.work>" RESEND_SIGNUP_TEMPLATE=signup
COPY --from=build /app/package.json ./
COPY --from=build /app/node_modules ./node_modules
COPY --from=build /app/dist ./dist
COPY migrations ./migrations
USER node
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://127.0.0.1:3000/health || exit 1
CMD ["node", "dist/index.js"]
