# Recode

Recode is a privacy-conscious, no-ad media utility platform for common image, video, and audio operations. It is designed to give users focused tools without requiring an account or installing desktop software.

The project is being built backend-first as a production-grade Go system, followed by a Next.js frontend.

## Planned tools

### Video

- Convert video to black and white
- Extract audio from video
- Remove audio and export a silent video
- Convert MP4, MOV, WebM, and MKV formats
- Clip video by start and end time

### Images

- Convert images to black and white
- Convert JPEG, PNG, WebP, and AVIF formats
- Compress images
- Resize images while preserving aspect ratio

## Architecture

```text
Next.js frontend
       |
       v
    Go API ------ PostgreSQL
       |              |
       v              |
     Redis queue <----|
       |
       v
   Go workers
       |
       +-- FFmpeg / ffprobe
       +-- libvips
       |
       v
Temporary media storage
```

The API accepts and validates requests, persists job metadata, and queues processing work. Separate workers execute CPU-intensive media operations so API availability is isolated from FFmpeg resource usage.

## Technology

- **Backend:** Go
- **Configuration:** Koanf
- **Database:** PostgreSQL
- **Queue and coordination:** Redis
- **Media processing:** FFmpeg, ffprobe, and libvips
- **Frontend:** Next.js, TypeScript, TanStack Query, and Zustand
- **Deployment:** Docker Compose and Caddy on a Hostinger VPS

## Engineering goals

Recode is also an interview-ready backend engineering project. Its implementation focuses on:

- Streaming uploads without loading complete files into memory
- Asynchronous jobs and explicit state transitions
- At-least-once delivery and idempotent processing
- Bounded concurrency and backpressure
- Child-process cancellation and timeouts
- Anonymous authorization and abuse controls
- Automatic retention and cleanup
- Structured logging, metrics, and operational health checks
- Integration, failure-path, race, and end-to-end testing
- Reproducible deployment, backup, recovery, and rollback

## Current status

Recode is under active development. The current backend foundation includes:

- A Go module with separate API and worker boundaries
- Layered configuration using defaults and `RECODE_*` environment variables
- Typed configuration decoding and validation
- Structured JSON logging
- An API composition root

The HTTP server, media job model, storage, queue, processors, and frontend are still being implemented.

## Repository layout

```text
recode/
├── backend/
│   ├── cmd/
│   │   ├── api/
│   │   └── worker/
│   └── internal/
│       ├── config/
│       └── observability/
└── README.md
```

## Development

Requirements:

- Go 1.26 or newer
- Docker with Docker Compose
- FFmpeg and libvips in the future worker environment

Run the current API bootstrap:

```bash
cd backend
go run ./cmd/api
```

Optional configuration overrides:

```bash
RECODE_HTTP_ADDRESS=127.0.0.1:9000 \
RECODE_SHUTDOWN_TIMEOUT=25s \
go run ./cmd/api
```

Run backend tests:

```bash
cd backend
go test ./...
```

## License

A license has not yet been selected.
