# Recode

Recode is an account-free, no-ad media processing API for common image and
video tasks. A Go API accepts uploads and returns immediately, while a separate
Go worker performs CPU-intensive work with FFmpeg.

This repository currently contains the minimum end-to-end backend. The frontend
will be built after this backend has been tested and revisited from first
principles.

## Supported operations

| Operation | `operation` value | Options |
| --- | --- | --- |
| Image to grayscale | `image_grayscale` | none |
| Convert image | `image_convert` | `format`: `jpg`, `jpeg`, `png`, or `webp` |
| Compress image | `image_compress` | `quality`: 1–100; default 80 |
| Resize image | `image_resize` | positive `width`, `height`, or both |
| Video to grayscale | `video_grayscale` | none |
| Extract audio | `video_extract_audio` | `format`: `mp3`, `wav`, or `m4a` |
| Remove audio | `video_remove_audio` | none |
| Convert video | `video_convert` | `format`: `mp4`, `webm`, or `mov` |
| Clip video | `video_clip` | non-negative `start_seconds` and positive `duration_seconds` |

Options are sent as one JSON object in the multipart `options` field.

## Architecture

```text
client
  │ multipart upload
  ▼
Go API ─────► local shared media volume
  │
  ├─────────► PostgreSQL (job state and metadata)
  │
  └─────────► Redis list (job ID)
                 │
                 ▼
              Go worker
                 ├── ffprobe validates the input
                 ├── FFmpeg processes it
                 ├── result goes to shared storage
                 └── final state goes to PostgreSQL
```

The API and worker are separate processes. A long conversion therefore does not
hold an HTTP request open or consume an API handler until it finishes.

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

## Run locally with Docker

Requirements:

- Docker with Docker Compose
- Go 1.26.4 or newer only when running tests outside containers

Start PostgreSQL and Redis:

```bash
docker compose up -d postgres redis
```

Apply every pending migration:

```bash
docker compose run --rm migrate
```

Build and start the API and worker:

```bash
docker compose build api worker
docker compose up -d api worker
```

Check dependency and API health:

```bash
docker compose ps
curl http://localhost:8080/health/live
curl http://localhost:8080/health/ready
```

The API is exposed on `localhost:8080`, PostgreSQL on `localhost:5433`, and
Redis on `localhost:6379`. Uploaded and generated files live in the Docker
`media_data` volume.

## Exercise the API

Create a resize job. Replace `/path/to/photo.png` with a real file:

```bash
curl -sS -X POST http://localhost:8080/api/v1/jobs \
  -F 'operation=image_resize' \
  -F 'options={"width":640}' \
  -F 'file=@/path/to/photo.png'
```

The `202 Accepted` response contains both an `id` and an `owner_token`:

```json
{
  "id": "example-job-id",
  "operation": "image_resize",
  "status": "queued",
  "progress": 0,
  "attempt": 0,
  "result_ready": false,
  "owner_token": "example-secret"
}
```

The token is shown only when the job is created. The future frontend must keep
it client-side and send it as a bearer token:

```bash
curl -sS http://localhost:8080/api/v1/jobs/JOB_ID \
  -H 'Authorization: Bearer OWNER_TOKEN'
```

Download a completed result:

```bash
curl -fLo result.jpg http://localhost:8080/api/v1/jobs/JOB_ID/result \
  -H 'Authorization: Bearer OWNER_TOKEN'
```

Cancel or delete:

```bash
curl -sS -X POST http://localhost:8080/api/v1/jobs/JOB_ID/cancel \
  -H 'Authorization: Bearer OWNER_TOKEN'

curl -i -X DELETE http://localhost:8080/api/v1/jobs/JOB_ID \
  -H 'Authorization: Bearer OWNER_TOKEN'
```

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

Repository tests are skipped unless a test database URL is supplied. After
running the migrations, include them with:

```bash
RECODE_TEST_DATABASE_URL='postgres://recode:recode_dev@localhost:5433/recode?sslmode=disable' \
go test ./internal/repository/postgres
```

## Configuration

Koanf loads defaults first and then overrides them from `RECODE_*` environment
variables. See [.env.example](.env.example) for the complete local set.

Important production settings include:

- `RECODE_DATABASE_URL`
- `RECODE_REDIS_ADDRESS`
- `RECODE_STORAGE_ROOT`
- `RECODE_UPLOAD_MAX_BYTES`
- `RECODE_RESULT_RETENTION`
- `RECODE_FFMPEG_PATH`
- `RECODE_FFPROBE_PATH`

## Current boundary

This is a deliberately minimal working backend, not yet a public-production
release. The ignored `defer.md` file is the engineering ledger for queue
durability, recovery, rate limits, process isolation, observability, TLS,
backups, and other hardening work that will be revisited after the first
end-to-end test.
