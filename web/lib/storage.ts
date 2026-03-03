import path from "path";

const DEFAULT_UPLOAD_DIR = path.resolve(process.cwd(), "../uploads");

export function getUploadDir() {
  return process.env.UPLOAD_DIR || DEFAULT_UPLOAD_DIR;
}

export function sanitizeFileName(name: string) {
  // Keep filenames filesystem-safe and avoid path traversal.
  return path
    .basename(name)
    .replace(/[^a-zA-Z0-9._-]/g, "_")
    .replace(/_+/g, "_");
}
