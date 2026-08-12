"use client";

import {
  type ChangeEvent,
  type DragEvent,
  useId,
  useRef,
  useState,
} from "react";
import {
  UploadSimple as FileUp,
  WarningCircle as CircleAlert,
} from "@phosphor-icons/react/ssr";

import type { ToolCategory } from "@/config/tools";
import { getAcceptAttribute } from "@/lib/media/file-rules";
import type { FileValidationFailure } from "@/lib/media/validate-file";

interface FileDropzoneProps {
  category: ToolCategory;
  error: FileValidationFailure | null;
  onFileSelected: (file: File | null) => boolean;
}

export function FileDropzone({
  category,
  error,
  onFileSelected,
}: FileDropzoneProps) {
  const inputId = useId();
  const helpId = `${inputId}-help`;
  const errorId = `${inputId}-error`;

  const inputRef = useRef<HTMLInputElement | null>(null);
  const dragDepthRef = useRef(0);

  const [dragActive, setDragActive] = useState(false);

  function handleInputChange(event: ChangeEvent<HTMLInputElement>) {
    const file = event.currentTarget.files?.item(0) ?? null;

    onFileSelected(file);

    event.currentTarget.value = "";
  }

  function handleDragEnter(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
    event.stopPropagation();

    dragDepthRef.current += 1;
    setDragActive(true);
  }

  function handleDragOver(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
    event.stopPropagation();
  }

  function handleDragLeave(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
    event.stopPropagation();

    dragDepthRef.current = Math.max(0, dragDepthRef.current - 1);

    if (dragDepthRef.current === 0) {
      setDragActive(false);
    }
  }

  function handleDrop(event: DragEvent<HTMLDivElement>) {
    event.preventDefault();
    event.stopPropagation();

    dragDepthRef.current = 0;
    setDragActive(false);

    const file = event.dataTransfer.files.item(0);
    onFileSelected(file);
  }

  function openFilePicker() {
    inputRef.current?.click();
  }

  return (
    <div>
      <div
        className={[
          "grid min-h-64 place-items-center",
          "rounded-md border border-dashed px-5 text-center",
          "transition-[border-color,background-color,transform]",
          "duration-200 ease-workshop",
          dragActive
            ? "border-ink bg-accent/40"
            : "border-line-strong bg-canvas hover:border-ink",
          error ? "border-danger" : "",
        ].join(" ")}
        onDragEnter={handleDragEnter}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        <input
          ref={inputRef}
          id={inputId}
          type="file"
          accept={getAcceptAttribute(category)}
          aria-describedby={error ? `${helpId} ${errorId}` : helpId}
          aria-invalid={Boolean(error)}
          className="sr-only"
          onChange={handleInputChange}
        />

        <div className="max-w-md">
          <span className="mx-auto grid size-10 place-items-center rounded-md bg-accent text-accent-strong">
            <FileUp aria-hidden="true" className="size-5" />
          </span>

          <h2 className="mt-4 text-base font-semibold tracking-[-0.02em]">
            {dragActive ? "Release to add the file" : "Drop one file here"}
          </h2>

          <p id={helpId} className="mt-2 text-xs leading-5 text-ink-muted">
            JPG, PNG or WebP · up to 1 GB
          </p>

          <button
            type="button"
            className="mt-5 min-h-10 rounded-md bg-ink px-4 py-2 text-xs font-semibold text-canvas shadow-hard hover:-translate-y-0.5 hover:shadow-soft"
            onClick={openFilePicker}
          >
            Browse files
          </button>

          {error && (
            <p
              id={errorId}
              role="alert"
              className="mt-6 flex items-center justify-center gap-2 text-sm font-semibold text-danger"
            >
              <CircleAlert aria-hidden="true" className="size-4 shrink-0" />
              {error.message}
            </p>
          )}
        </div>
      </div>
    </div>
  );
}
