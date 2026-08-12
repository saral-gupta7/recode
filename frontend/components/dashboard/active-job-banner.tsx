"use client";

import Link from "next/link";
import { ArrowRight, SpinnerGap as LoaderCircle } from "@phosphor-icons/react/ssr";

import { getToolBySlug } from "@/config/tools";
import { useWorkspaceStore } from "@/stores/use-workspace-store";

export function ActiveJobBanner() {
  const activeJob = useWorkspaceStore((state) => state.activeJob);
  const hydrated = useWorkspaceStore((state) => state.hydrated);

  if (!hydrated || !activeJob) {
    return null;
  }

  const tool = getToolBySlug(activeJob.toolSlug);

  return (
    <aside
      aria-label="Active job"
      className="mt-5 flex flex-col gap-3 rounded-md border border-line bg-accent p-3 sm:flex-row sm:items-center sm:justify-between"
    >
      <div className="flex items-center gap-3">
        <LoaderCircle
          aria-hidden="true"
          className="size-4 shrink-0 animate-spin text-accent-strong"
        />
        <p className="text-xs">
          <strong className="font-bold">You have an active job.</strong>{" "}
          <span className="text-ink-muted">
            Continue {tool?.title ?? "processing"} to view its latest state.
          </span>
        </p>
      </div>
      <Link
        href={`/tools/${activeJob.toolSlug}`}
        className="inline-flex min-h-9 shrink-0 items-center justify-center gap-2 rounded-md bg-ink px-3 py-1.5 text-xs font-semibold text-canvas hover:opacity-90"
      >
        Open workspace
        <ArrowRight aria-hidden="true" className="size-4" />
      </Link>
    </aside>
  );
}
