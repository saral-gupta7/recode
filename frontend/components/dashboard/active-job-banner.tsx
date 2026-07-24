"use client";

import Link from "next/link";
import { ArrowRight, LoaderCircle } from "lucide-react";

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
      className="mt-6 flex flex-col gap-4 rounded-lg border border-line bg-surface p-4 sm:flex-row sm:items-center sm:justify-between"
    >
      <div className="flex items-center gap-3">
        <LoaderCircle
          aria-hidden="true"
          className="size-5 shrink-0 animate-spin text-accent-strong"
        />
        <p className="text-sm">
          <strong className="font-bold">You have an active job.</strong>{" "}
          <span className="text-ink-muted">
            Continue {tool?.title ?? "processing"} to view its latest state.
          </span>
        </p>
      </div>
      <Link
        href={`/tools/${activeJob.toolSlug}`}
        className="inline-flex min-h-11 shrink-0 items-center justify-center gap-2 rounded-md bg-ink px-4 py-2 text-sm font-semibold text-surface hover:opacity-85"
      >
        Open workspace
        <ArrowRight aria-hidden="true" className="size-4" />
      </Link>
    </aside>
  );
}
