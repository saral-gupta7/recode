"use client";

import { useEffect, useState } from "react";
import {
  ArrowClockwise as RotateCcw,
  DownloadSimple as Download,
} from "@phosphor-icons/react/ssr";

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

  return (
    <section className="overflow-hidden rounded-lg border border-line bg-surface shadow-hard">
      <div className="grid min-h-72 place-items-center overflow-hidden bg-inverse p-4">
        {(
          // The source is an authenticated browser-local Blob URL.
          // eslint-disable-next-line @next/next/no-img-element
          <img
            src={objectURL}
            alt={`Processed result for ${tool.title}`}
            className="max-h-[30rem] w-full object-contain"
          />
        )}

      </div>

      <div className="grid gap-2 border-t border-line p-3 sm:grid-cols-2">
        <a
          href={objectURL}
          download={result.filename}
          className="inline-flex min-h-10 items-center justify-center gap-2 rounded-md bg-ink px-4 py-2 text-xs font-semibold text-canvas"
        >
          <Download aria-hidden="true" className="size-4" />
          Download result
        </a>
        <button
          type="button"
          disabled={deleting}
          className="inline-flex min-h-10 items-center justify-center gap-2 rounded-md border border-line bg-surface px-4 py-2 text-xs font-semibold hover:bg-canvas disabled:opacity-50"
          onClick={onStartAnother}
        >
          <RotateCcw aria-hidden="true" className="size-4" />
          {deleting ? "Deleting…" : "Start another"}
        </button>
      </div>
    </section>
  );
}
