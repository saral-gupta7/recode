"use client";

import { CircleAlert, RotateCcw } from "lucide-react";

export default function ErrorPage({
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <main
      id="main-content"
      className="grid min-h-[70svh] place-items-center px-6 py-16"
    >
      <section className="w-full max-w-2xl border-2 border-danger bg-surface p-8 text-center sm:p-12">
        <CircleAlert
          aria-hidden="true"
          className="mx-auto size-10 text-danger"
        />
        <p className="mt-6 font-mono text-xs font-semibold uppercase tracking-[0.14em] text-danger">
          Interface fault
        </p>
        <h1 className="mt-4 text-4xl font-black uppercase tracking-[-0.05em]">
          Something went wrong
        </h1>
        <p className="mx-auto mt-4 max-w-md leading-7 text-ink-muted">
          Your active job credentials remain in this browser session. Retry the
          interface without creating a duplicate job.
        </p>
        <button
          type="button"
          className="mt-8 inline-flex min-h-12 items-center gap-2 border-2 border-ink bg-accent px-5 py-3 text-sm font-bold text-[#171a17] shadow-hard"
          onClick={reset}
        >
          <RotateCcw aria-hidden="true" className="size-4" />
          Try again
        </button>
      </section>
    </main>
  );
}
