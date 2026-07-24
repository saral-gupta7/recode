"use client";

import { useState } from "react";
import { FileWarning, Trash2 } from "lucide-react";

import type { ToolCategory } from "@/config/tools";
import type { SelectedFile } from "@/hooks/use-selected-file";
import { formatFileSize } from "@/lib/media/format-file-size";

interface FilePreviewProps {
  category: ToolCategory;
  selectedFile: SelectedFile;
  onRemove: () => void;
}

export function FilePreview({
  category,
  selectedFile,
  onRemove,
}: FilePreviewProps) {
  const [previewFailed, setPreviewFailed] = useState(false);
  const { file, previewUrl } = selectedFile;

  return (
    <section className="overflow-hidden rounded-lg border border-line bg-surface">
      <div className="relative grid min-h-[25rem] place-items-center overflow-hidden bg-inverse">
        {!previewFailed && category === "image" && (
          // Blob URLs are browser-local and cannot use Next.js optimization.
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={previewUrl}
            alt={`Preview of ${file.name}`}
            className="max-h-[34rem] w-full object-contain"
            onError={() => setPreviewFailed(true)}
          />
        )}

        {!previewFailed && category === "video" && (
          <video
            src={previewUrl}
            controls
            preload="metadata"
            aria-label={`Preview of ${file.name}`}
            className="max-h-[34rem] w-full bg-black object-contain"
            onError={() => setPreviewFailed(true)}
          />
        )}

        {previewFailed && (
          <div
            role="alert"
            className="max-w-md px-6 py-16 text-center text-on-inverse"
          >
            <FileWarning
              aria-hidden="true"
              className="mx-auto size-10 text-warning"
            />
            <h2 className="mt-5 text-xl font-bold">Preview unavailable</h2>
            <p className="mt-3 text-sm leading-6 text-on-inverse/65">
              The browser could not display this file. FFmpeg may still support
              it, or you can remove it and choose another.
            </p>
          </div>
        )}
      </div>

      <div className="grid gap-4 border-t border-line p-4 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center sm:p-5">
        <div className="min-w-0">
          <p title={file.name} className="truncate text-sm font-bold">
            {file.name}
          </p>
          <div className="mt-2 flex flex-wrap gap-x-3 gap-y-1 font-mono text-[0.68rem] font-semibold uppercase tracking-[0.08em] text-ink-muted">
            <span>{formatFileSize(file.size)}</span>
            <span aria-hidden="true">/</span>
            <span>{file.type || "Unknown type"}</span>
          </div>
        </div>

        <button
          type="button"
          className="inline-flex min-h-11 items-center justify-center gap-2 rounded-md border border-line bg-surface px-4 py-2 text-sm font-semibold hover:bg-canvas"
          onClick={onRemove}
        >
          <Trash2 aria-hidden="true" className="size-4" />
          Remove file
        </button>
      </div>
    </section>
  );
}
