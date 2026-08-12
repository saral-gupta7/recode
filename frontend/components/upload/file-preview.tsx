"use client";

import { useState, type CSSProperties } from "react";
import {
  ImageBroken as FileWarning,
  Trash as Trash2,
} from "@phosphor-icons/react/ssr";

import type { SelectedFile } from "@/hooks/use-selected-file";
import { formatFileSize } from "@/lib/media/format-file-size";

interface FilePreviewProps {
  selectedFile: SelectedFile;
  imageStyle?: CSSProperties;
  onRemove: () => void;
}

export function FilePreview({
  selectedFile,
  onRemove,
  imageStyle,
}: FilePreviewProps) {
  const [previewFailed, setPreviewFailed] = useState(false);
  const { file, previewUrl } = selectedFile;

  return (
    <section className="overflow-hidden rounded-md border border-line bg-surface">
      <div className="relative grid min-h-64 place-items-center overflow-hidden bg-inverse p-3">
        {!previewFailed && (
          // Blob URLs are browser-local and cannot use Next.js optimization.
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={previewUrl}
            alt={`Preview of ${file.name}`}
            style={imageStyle}
            className="max-h-[26rem] w-full rounded-sm object-contain"
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
              The browser could not display this image. Remove it and choose
              another supported file.
            </p>
          </div>
        )}
      </div>

      <div className="grid gap-3 border-t border-line p-3 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center">
        <div className="min-w-0">
          <p title={file.name} className="truncate text-xs font-semibold">
            {file.name}
          </p>
          <div className="mt-1 flex flex-wrap gap-x-2 gap-y-1 font-mono text-[0.6rem] text-ink-muted">
            <span>{formatFileSize(file.size)}</span>
            <span aria-hidden="true">/</span>
            <span>{file.type || "Unknown type"}</span>
          </div>
        </div>

        <button
          type="button"
          className="inline-flex min-h-9 items-center justify-center gap-2 rounded-md border border-line bg-surface px-3 py-1.5 text-xs font-semibold hover:bg-canvas"
          onClick={onRemove}
        >
          <Trash2 aria-hidden="true" className="size-4" />
          Remove file
        </button>
      </div>
    </section>
  );
}
