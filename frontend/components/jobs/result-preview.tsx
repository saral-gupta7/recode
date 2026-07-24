"use client";

import { useEffect, useState } from "react";
import { Download, RotateCcw } from "lucide-react";

import type { MediaTool } from "@/config/tools";
import type { ResultFile } from "@/lib/api/jobs";

interface ResultPreviewProps {
  tool: MediaTool;
  result: ResultFile;
  deleting: boolean;
  onStartAnother: () => void;
}

export function ResultPreview({
  tool,
  result,
  deleting,
  onStartAnother,
}: ResultPreviewProps) {
  const [objectURL] = useState(() => URL.createObjectURL(result.blob));

  useEffect(() => {
    return () => URL.revokeObjectURL(objectURL);
  }, [objectURL]);

  const isAudio = tool.operation === "video_extract_audio";

  return (
    <section className="overflow-hidden rounded-lg border border-line bg-surface">
      <div className="grid min-h-[27rem] place-items-center overflow-hidden bg-inverse p-4">
        {tool.category === "image" && (
          // The source is an authenticated browser-local Blob URL.
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={objectURL}
            alt={`Processed result for ${tool.title}`}
            className="max-h-[38rem] w-full object-contain"
          />
        )}

        {tool.category === "video" && !isAudio && (
          <video
            src={objectURL}
            controls
            autoPlay={false}
            className="max-h-[38rem] w-full bg-black object-contain"
          />
        )}

        {isAudio && (
          <div className="w-full max-w-xl border-2 border-on-inverse/25 bg-on-inverse/5 p-6 sm:p-10">
            <p className="font-mono text-[0.7rem] font-semibold uppercase tracking-[0.14em] text-accent">
              Audio result
            </p>
            <audio src={objectURL} controls className="mt-8 w-full" />
          </div>
        )}
      </div>

      <div className="grid gap-3 border-t border-line p-4 sm:grid-cols-2 sm:p-5">
        <a
          href={objectURL}
          download={result.filename}
          className="inline-flex min-h-12 items-center justify-center gap-2 rounded-md bg-ink px-5 py-3 text-sm font-semibold text-surface"
        >
          <Download aria-hidden="true" className="size-4" />
          Download result
        </a>
        <button
          type="button"
          disabled={deleting}
          className="inline-flex min-h-12 items-center justify-center gap-2 rounded-md border border-line bg-surface px-5 py-3 text-sm font-semibold hover:bg-canvas disabled:opacity-50"
          onClick={onStartAnother}
        >
          <RotateCcw aria-hidden="true" className="size-4" />
          {deleting ? "Deleting…" : "Start another"}
        </button>
      </div>
    </section>
  );
}
