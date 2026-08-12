"use client";

import {
  ArrowClockwise as RotateCcw,
  WarningCircle as CircleAlert,
} from "@phosphor-icons/react/ssr";

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
      <section className="w-full max-w-xl rounded-lg border border-line bg-surface p-8 text-center shadow-hard sm:p-10">
        <CircleAlert
          aria-hidden="true"
          className="mx-auto size-10 text-danger"
        />
        <p className="mt-5 text-xs font-semibold text-danger">
          Interface error
        </p>
        <h1 className="mt-3 text-3xl font-semibold tracking-[-0.04em]">
          Something went wrong
        </h1>
        <p className="mx-auto mt-4 max-w-md leading-7 text-ink-muted">
          Your active job credentials remain in this browser session. Retry the
          interface without creating a duplicate job.
        </p>
        <button
          type="button"
          className="mt-7 inline-flex min-h-10 items-center gap-2 rounded-md bg-ink px-4 py-2 text-sm font-semibold text-canvas"
          onClick={reset}
        >
          <RotateCcw aria-hidden="true" className="size-4" />
          Try again
        </button>
      </section>
    </main>
  );
}
