"use client";

import {
  ArrowClockwise as RotateCcw,
  Clock as Clock3,
  Prohibit as Ban,
  SpinnerGap as LoaderCircle,
  WarningCircle as CircleAlert,
  X,
} from "@phosphor-icons/react/ssr";

import type { Job } from "@/lib/api/contracts";

interface JobStatusPanelProps {
  job: Job;
  cancelling: boolean;
  deleting: boolean;
  onCancel: () => void;
  onReset: (keepFile: boolean) => void;
}

const failureMessages: Record<string, string> = {
  invalid_media: "The uploaded file could not be read as supported media.",
  unsupported_media: "This file does not support the selected operation.",
  processing_failed: "FFmpeg could not complete this transformation.",
  input_read_failed: "The worker could not read the uploaded file.",
  output_store_failed: "The result could not be stored.",
};

export function JobStatusPanel({
  job,
  cancelling,
  deleting,
  onCancel,
  onReset,
}: JobStatusPanelProps) {
  if (job.status === "queued" || job.status === "processing") {
    return (
      <StatusFrame
        icon={<LoaderCircle className="size-10 animate-spin text-accent" />}
        eyebrow={job.status === "queued" ? "Job queued" : "Recode is working"}
        title={
          job.status === "queued"
            ? "Waiting for the worker"
            : "Processing your file"
        }
        description={
          job.status === "queued"
            ? "The API accepted the job. It will begin as soon as the worker is available."
            : "This can take a while for large or high-resolution media. You can keep this tab open."
        }
      >
        <button
          type="button"
          disabled={cancelling}
          className="inline-flex min-h-11 items-center gap-2 rounded-md border border-line px-4 py-2 text-sm font-semibold hover:bg-canvas disabled:opacity-50"
          onClick={onCancel}
        >
          <X aria-hidden="true" className="size-4" />
          {cancelling ? "Requesting cancellation…" : "Cancel job"}
        </button>
      </StatusFrame>
    );
  }

  if (job.status === "cancelling") {
    return (
      <StatusFrame
        icon={<Clock3 className="size-10 text-warning" />}
        eyebrow="Cancellation requested"
        title="Stopping safely"
        description="The worker has received the cancellation request. It may need a moment to stop the active media process."
      />
    );
  }

  if (job.status === "failed") {
    return (
      <StatusFrame
        icon={<CircleAlert className="size-10 text-danger" />}
        eyebrow="Processing failed"
        title="This job did not finish"
        description={
          failureMessages[job.failure_code ?? ""] ??
          "The media processor could not complete this job."
        }
      >
        <button
          type="button"
          disabled={deleting}
          className="inline-flex min-h-11 items-center gap-2 rounded-md bg-ink px-4 py-2 text-sm font-semibold text-surface disabled:opacity-50"
          onClick={() => onReset(true)}
        >
          <RotateCcw aria-hidden="true" className="size-4" />
          {deleting ? "Cleaning up…" : "Try again"}
        </button>
      </StatusFrame>
    );
  }

  if (job.status === "cancelled") {
    return (
      <StatusFrame
        icon={<Ban className="size-10 text-warning" />}
        eyebrow="Job cancelled"
        title="Processing stopped"
        description="The job is no longer running. Its temporary data can now be removed."
      >
        <button
          type="button"
          disabled={deleting}
          className="min-h-11 rounded-md bg-ink px-4 py-2 text-sm font-semibold text-surface disabled:opacity-50"
          onClick={() => onReset(false)}
        >
          {deleting ? "Cleaning up…" : "Start another job"}
        </button>
      </StatusFrame>
    );
  }

  if (job.status === "expired") {
    return (
      <StatusFrame
        icon={<Clock3 className="size-10 text-warning" />}
        eyebrow="Result expired"
        title="This file is no longer available"
        description="Recode automatically removed the retained media. Clear this session to begin again."
      >
        <button
          type="button"
          className="min-h-11 rounded-md bg-ink px-4 py-2 text-sm font-semibold text-surface"
          onClick={() => onReset(false)}
        >
          Start another job
        </button>
      </StatusFrame>
    );
  }

  return null;
}

function StatusFrame({
  icon,
  eyebrow,
  title,
  description,
  children,
}: {
  icon: React.ReactNode;
  eyebrow: string;
  title: string;
  description: string;
  children?: React.ReactNode;
}) {
  return (
    <section
      aria-live="polite"
      className="grid min-h-80 place-items-center rounded-lg border border-line bg-surface px-6 py-10 text-center shadow-hard"
    >
      <div className="max-w-xl">
        <div className="mx-auto grid size-14 place-items-center rounded-lg bg-canvas">
          {icon}
        </div>
        <p className="mt-4 text-xs font-medium text-accent-strong">
          {eyebrow}
        </p>
        <h2 className="mt-2 text-xl font-semibold tracking-[-0.035em] sm:text-2xl">
          {title}
        </h2>
        <p className="mt-3 text-sm leading-6 text-ink-muted">
          {description}
        </p>
        {children && <div className="mt-6">{children}</div>}
      </div>
    </section>
  );
}
