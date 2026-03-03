# Recode

Recode is a login-free video-to-GIF conversion app with an asynchronous processing pipeline.

You upload a video in the web UI, the backend enqueues a job in RabbitMQ, a worker converts the video with FFmpeg, and the UI polls until the GIF is ready to download.

## Features

- Drag-and-drop video upload (`react-dropzone`)
- Upload state UX (`idle` -> `uploading` -> `processing` -> `done`)
- Animated UI transitions with Framer Motion
- Async job queue using RabbitMQ (`video_processing` queue)
- Background media processing worker with `fluent-ffmpeg`
- Polling-based processing status endpoint
- Direct download endpoint for completed GIF files
- Shared local storage for uploaded and processed files (`/uploads`)
- Durable RabbitMQ messages and explicit ack/reject handling in worker

## Tech Stack

### Web (`/web`)

- Next.js 16 (App Router)
- React 19 + TypeScript
- Tailwind CSS 4
- `react-dropzone` for file intake
- `motion` + `lucide-react` for UI animation and icons
- `amqplib` for RabbitMQ publishing

### Worker (`/worker`)

- Bun runtime + TypeScript
- `amqplib` for consuming queue messages
- `fluent-ffmpeg` for video conversion

### Infrastructure

- RabbitMQ 3 management image via Docker Compose
- Local filesystem storage for media files

## Repository Structure

```text
recode/
├─ docker-compose.yml        # RabbitMQ service
├─ uploads/                  # Uploaded videos + generated GIFs
├─ web/                      # Next.js frontend + API routes
│  ├─ app/page.tsx           # Drag/drop UI + polling logic
│  └─ app/api/
│     ├─ upload/route.ts     # Saves upload + enqueues job
│     ├─ status/route.ts     # Checks processed file existence
│     └─ download/route.ts   # Serves generated GIF as attachment
└─ worker/
   └─ src/index.ts           # RabbitMQ consumer + FFmpeg conversion
```

## How It Works

1. User drops a video in the UI (`web/app/page.tsx`).
2. Frontend posts multipart form data to `POST /api/upload`.
3. API route saves file to `../uploads` and enqueues job metadata to `video_processing`.
4. Worker consumes the job, runs FFmpeg, and writes:
   - `processed-<original-upload-name>.gif`
5. Frontend polls `GET /api/status?fileName=<upload-name>` every 2 seconds.
6. When the output exists, API returns a download URL.
7. User downloads through `GET /api/download?file=<processed-name>`.

## API Endpoints

### `POST /api/upload`

Accepts multipart form data:

- `video` (required): video file
- `title` (optional): currently passed through in response
- `description` (optional): currently passed through in response

Behavior:

- Saves file into `uploads/` using a timestamp-prefixed name
- Publishes conversion job to RabbitMQ queue `video_processing`
- Returns JSON including `fileName` used for status checks

### `GET /api/status?fileName=<name>`

Behavior:

- Checks whether `uploads/processed-<fileName>.gif` exists
- Returns:
  - `{ "status": "processing" }` when not ready
  - `{ "status": "done", "downloadUrl": "/api/download?file=..." }` when ready

### `GET /api/download?file=<processed-file-name>`

Behavior:

- Reads requested file from `uploads/`
- Returns attachment response with `Content-Type: image/gif`

## Local Development

## Prerequisites

- Bun (recommended for this repo)
- Docker + Docker Compose
- FFmpeg installed on your machine and available in `PATH`

## 1) Start RabbitMQ

From repo root:

```bash
docker compose up -d
```

RabbitMQ services:

- AMQP: `localhost:5672`
- Management UI: `http://localhost:15672` (guest / guest)

## 2) Install dependencies

```bash
cd web && bun install
cd ../worker && bun install
```

## 3) Run worker

In one terminal:

```bash
cd worker
bun run src/index.ts
```

## 4) Run web app

In another terminal:

```bash
cd web
bun run dev
```

Open `http://localhost:3000`.

## Production Notes

- Current RabbitMQ connection strings are hardcoded to `amqp://localhost:5672`.
- File storage is local disk only (`uploads/`); no object storage integration yet.
- Status tracking is file-existence based; no DB/state store currently used.
- Queue and output naming conventions are fixed in code.

## Known Limitations (Current Implementation)

- No authentication or authorization.
- No upload size/type validation beyond client-side accepted MIME group.
- No cleanup policy for old uploads/outputs.
- No retry/backoff strategy beyond queue reject behavior.
- No progress percentages (only processing complete vs not complete).
- Some scaffolding files remain placeholders (`web/lib/redis.ts`, `web/lib/rabbitmq.ts`, empty `web/Dockerfile`).

## Future Improvements

- Move broker URL, queue name, and paths to environment variables
- Add file validation and max size enforcement
- Add signed/expiring download links and basic abuse controls
- Add persistent job state (Redis/Postgres) for richer statuses
- Add output options (fps, duration, resolution, crop)
- Add automatic cleanup/retention jobs
- Containerize worker + web for one-command full stack startup

## License

No license file is currently defined in this repository.
