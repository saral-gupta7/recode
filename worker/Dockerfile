FROM oven/bun:1.1-alpine AS builder
WORKDIR /app

COPY package.json bun.lockb ./
RUN bun install --frozen-lockfile

COPY . .
RUN bun run build 

FROM node:18-bookworm-slim AS runner
WORKDIR /app

RUN apt-get update && apt-get install -y ffmpeg && rm -rf /var/lib/apt/lists/*

ENV NODE_ENV=production

COPY package.json bun.lockb ./
RUN npm install -g bun && bun install --production

COPY --from=builder /app/dist ./dist

RUN mkdir -p /app/uploads

CMD ["node", "dist/index.js"]