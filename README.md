# Recode

Recode is an account-free, ad-free web application for focused image-processing
tasks. A Next.js interface provides local previews and job progress, a Go API
accepts uploads without holding long-running requests open, and a separate Go
worker performs CPU-intensive FFmpeg processing.

## Supported operations

| Operation | `operation` value | Options |
| --- | --- | --- |
| Image to grayscale | `image_grayscale` | none |
| Convert image | `image_convert` | `format`: `jpg`, `jpeg`, `png`, or `webp` |
| Compress image | `image_compress` | `quality`: 1–100; default 80 |
| Resize image | `image_resize` | positive `width`, `height`, or both |
| Crop image | `image_crop` | `x`, `y`, `width`, and `height` |
| Rotate image | `image_rotate` | `angle`: `90`, `180`, or `270` |
| Flip image | `image_flip` | `flip_direction`: `horizontal` or `vertical` |
| Generate thumbnail | `image_thumbnail` | `preset`: `square`, `preview`, or `social` |
| Strip metadata | `image_strip_metadata` | none |
| Adjust colour | `image_adjust` | `brightness`, `contrast`, and `saturation` |
| Blur image | `image_blur` | `strength`: 0.1–20 |
| Sharpen image | `image_sharpen` | `strength`: 0.1–5 |
| Pixelate image | `image_pixelate` | `block_size`: 2–100 |
| Add canvas padding | `image_padding` | per-side padding and `background` colour |

Options are sent as one JSON object in the multipart `options` field.

## Architecture

```text
client
  │ multipart upload
  ▼
Go API ─────► shared media volume
  │
  ├─────────► Neon PostgreSQL (job state and metadata)
  │
  └─────────► Redis list (job ID)
                 │
                 ▼
              Go worker
                 ├── ffprobe validates the input
                 ├── FFmpeg processes it
                 ├── result goes to shared storage
              └── final state goes to Neon PostgreSQL
```

The API and worker are separate processes. A long image transformation does not
hold an HTTP request open or consume an API handler until it finishes.

In production the browser talks only to Next.js. Same-origin rewrites forward
`/api/*` and `/health/*` to the private Go API origin, avoiding browser CORS
configuration and keeping internal service addresses out of the client bundle.

## Job lifecycle

```text
queued ──► processing ──► completed ──► expired
  │              │             cleanup removes files
  │              ├──► failed ────────► expired
  │              └──► cancelling ─► cancelled ─► expired
  └──────────────────────────────► cancelled
```

Completed, failed, and cancelled jobs expire after the configured retention
period. The worker performs periodic cleanup.

## Backend layout

```text
backend/
├── cmd/
│   ├── api/                 # API composition root
│   └── worker/              # worker composition root
├── internal/
│   ├── application/jobs/    # use cases and anonymous authorization
│   ├── config/              # Koanf defaults, environment loading, validation
│   ├── database/            # PostgreSQL pool
│   ├── httpapi/             # Gin routes, handlers, middleware, responses
│   ├── job/                 # job aggregate and state machine
│   ├── media/               # FFmpeg and ffprobe adapters
│   ├── observability/       # structured logging and request context
│   ├── queue/               # queue abstraction and Redis adapter
│   ├── repository/postgres/ # PostgreSQL repositories
│   ├── storage/             # bounded local object storage
│   ├── task/                # execution metadata and operation options
│   └── worker/              # dequeue, process, persist, and cleanup flow
└── migrations/
```

## Frontend layout

```text
frontend/
├── app/                     # Next.js routes and application boundaries
├── components/
│   ├── dashboard/           # catalog and active-job entry points
│   ├── jobs/                # lifecycle and result states
│   ├── layout/              # responsive navigation
│   ├── providers/           # theme and query clients
│   ├── tools/               # operation workspace and controls
│   └── upload/              # selection, validation, and previews
├── config/                  # typed tool registry
├── features/media-jobs/     # operation defaults and validation
├── hooks/                   # browser file/object URL lifecycle
├── lib/
│   ├── api/                 # transport, contracts, and runtime parsing
│   └── media/               # file rules and formatting
└── stores/                  # minimal persisted anonymous job credentials
```

TanStack Query owns server state: job creation, polling, cancellation, deletion,
and result retrieval. Zustand stores only cross-route client state. The anonymous
owner token is persisted in `sessionStorage`, not local storage, and is never
logged or encoded into a URL.

## Deployment configuration

The production stack uses Neon PostgreSQL, a local Redis queue, and a shared
Docker volume for uploaded and processed images. Copy [.env.example](.env.example)
to an ignored repository-root `.env` and replace every placeholder with the
direct, non-pooled connection string from Neon's **Connect** dialog.

```dotenv
RECODE_COMPOSE_DATABASE_URL=postgresql://USER:PASSWORD@NEON_HOST/neondb?sslmode=require&channel_binding=require
```

Never commit `.env` or paste its connection string into logs, documentation, or
issues. Compose requires this value and refuses to start when it is missing. The
one-shot migration service uses the same direct connection and must complete
successfully before the API and worker start. Redis and generated media remain
on persistent Docker volumes.

## API summary

| Method | Path | Success |
| --- | --- | --- |
| `POST` | `/api/v1/jobs` | `202`; creates and queues a job |
| `GET` | `/api/v1/jobs/:jobID` | `200`; returns current state |
| `POST` | `/api/v1/jobs/:jobID/cancel` | `200`; requests cancellation |
| `GET` | `/api/v1/jobs/:jobID/result` | `200`; streams completed result |
| `DELETE` | `/api/v1/jobs/:jobID` | `204`; removes a non-running job |
| `GET` | `/health/live` | process liveness |
| `GET` | `/health/ready` | PostgreSQL and Redis readiness |

All job endpoints except creation require
`Authorization: Bearer <owner_token>`.

## Run tests

Install/synchronize module dependencies and format the Go source:

```bash
cd backend
go mod tidy
go fmt ./...
```

Run unit tests:

```bash
go test ./...
```

Repository tests are skipped unless `RECODE_TEST_DATABASE_URL` points to an
isolated, migrated test database. Never run repository tests against production.

Validate the frontend:

```bash
cd frontend
pnpm lint
pnpm typecheck
pnpm test:run
pnpm build
```

## Configuration

Koanf loads application defaults first and then applies `RECODE_*` environment
overrides. Docker Compose maps the required `RECODE_COMPOSE_DATABASE_URL` secret
to `RECODE_DATABASE_URL` inside the API and worker containers.

Important production settings include:

- `RECODE_COMPOSE_DATABASE_URL` (Compose host environment)
- `RECODE_DATABASE_URL` (API and worker containers)
- `RECODE_REDIS_ADDRESS`
- `RECODE_STORAGE_ROOT`
- `RECODE_UPLOAD_MAX_BYTES`
- `RECODE_RESULT_RETENTION`
- `RECODE_FFMPEG_PATH`
- `RECODE_FFPROBE_PATH`

## Technology

- Frontend: Next.js, React, TypeScript, Tailwind CSS, TanStack Query, Zustand,
  Zod, next-themes, Vitest
- Backend: Go, Gin, Koanf, PostgreSQL, Redis, FFmpeg/ffprobe
- Delivery: Docker multi-stage builds and Docker Compose

The ignored `defer.md` file remains the engineering ledger for public-launch
hardening such as queue recovery, rate limits, process isolation, observability,
TLS, backups, and abuse controls.
