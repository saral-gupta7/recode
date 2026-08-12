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

## Technology

- Frontend: Next.js, React, TypeScript, Tailwind CSS, TanStack Query, Zustand,
  Zod, next-themes, and Vitest
- Backend: Go, Gin, Koanf, PostgreSQL, Redis, and FFmpeg/ffprobe
- Delivery: Docker multi-stage builds and Docker Compose
