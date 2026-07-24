"use client";

import type { JobOperation } from "@/config/tools";
import {
  type JobOptions,
  type OptionErrors,
} from "@/features/media-jobs/options";

interface ToolOptionsProps {
  operation: JobOperation;
  options: JobOptions;
  errors: OptionErrors;
  onChange: (options: JobOptions) => void;
}

const fieldStyles =
  "mt-2 min-h-12 w-full rounded-md border-2 border-line-strong bg-canvas px-3 text-sm outline-none transition-colors focus:border-ink";
const labelStyles =
  "block font-mono text-[0.68rem] font-semibold uppercase tracking-[0.12em] text-ink-muted";

export function ToolOptions({
  operation,
  options,
  errors,
  onChange,
}: ToolOptionsProps) {
  if (operation === "image_convert") {
    return (
      <SelectField
        label="Output format"
        value={options.format ?? "png"}
        values={["png", "jpg", "webp"]}
        onChange={(format) => onChange({ ...options, format })}
      />
    );
  }

  if (operation === "video_extract_audio") {
    return (
      <SelectField
        label="Audio format"
        value={options.format ?? "mp3"}
        values={["mp3", "wav", "m4a"]}
        onChange={(format) => onChange({ ...options, format })}
      />
    );
  }

  if (operation === "video_convert") {
    return (
      <SelectField
        label="Output format"
        value={options.format ?? "mp4"}
        values={["mp4", "webm", "mov"]}
        onChange={(format) => onChange({ ...options, format })}
      />
    );
  }

  if (operation === "image_compress") {
    const quality = options.quality ?? 80;

    return (
      <div>
        <div className="flex items-center justify-between">
          <label htmlFor="quality" className={labelStyles}>
            Output quality
          </label>
          <output
            htmlFor="quality"
            className="font-mono text-sm font-semibold"
          >
            {quality}%
          </output>
        </div>
        <input
          id="quality"
          type="range"
          min="1"
          max="100"
          value={quality}
          className="mt-4 w-full accent-[var(--color-accent-strong)]"
          onChange={(event) =>
            onChange({ ...options, quality: Number(event.target.value) })
          }
        />
        <FieldError message={errors.quality} />
      </div>
    );
  }

  if (operation === "image_resize") {
    return (
      <div className="grid gap-4 sm:grid-cols-2">
        <NumberField
          id="width"
          label="Width"
          suffix="PX"
          value={options.width}
          error={errors.width}
          onChange={(width) => onChange({ ...options, width })}
        />
        <NumberField
          id="height"
          label="Height"
          suffix="PX"
          value={options.height}
          error={errors.height}
          onChange={(height) => onChange({ ...options, height })}
        />
        <p className="sm:col-span-2 text-xs leading-5 text-ink-muted">
          Leave one dimension empty to preserve the original aspect ratio.
        </p>
      </div>
    );
  }

  if (operation === "video_clip") {
    return (
      <div className="grid gap-4 sm:grid-cols-2">
        <NumberField
          id="start-seconds"
          label="Start time"
          suffix="SEC"
          min={0}
          step={0.1}
          value={options.start_seconds}
          error={errors.start_seconds}
          onChange={(start_seconds) =>
            onChange({ ...options, start_seconds })
          }
        />
        <NumberField
          id="duration-seconds"
          label="Duration"
          suffix="SEC"
          min={0.1}
          step={0.1}
          value={options.duration_seconds}
          error={errors.duration_seconds}
          onChange={(duration_seconds) =>
            onChange({ ...options, duration_seconds })
          }
        />
      </div>
    );
  }

  return (
    <p className="text-sm leading-6 text-ink-muted">
      This operation has no additional settings. Recode will use a balanced
      output profile.
    </p>
  );
}

interface SelectFieldProps {
  label: string;
  value: string;
  values: readonly string[];
  onChange: (value: string) => void;
}

function SelectField({
  label,
  value,
  values,
  onChange,
}: SelectFieldProps) {
  return (
    <label className={labelStyles}>
      {label}
      <select
        value={value}
        className={fieldStyles}
        onChange={(event) => onChange(event.target.value)}
      >
        {values.map((item) => (
          <option key={item} value={item}>
            {item.toUpperCase()}
          </option>
        ))}
      </select>
    </label>
  );
}

interface NumberFieldProps {
  id: string;
  label: string;
  suffix: string;
  value: number | undefined;
  error?: string;
  min?: number;
  step?: number;
  onChange: (value: number | undefined) => void;
}

function NumberField({
  id,
  label,
  suffix,
  value,
  error,
  min = 1,
  step = 1,
  onChange,
}: NumberFieldProps) {
  const errorID = `${id}-error`;

  return (
    <div>
      <label htmlFor={id} className={labelStyles}>
        {label}
      </label>
      <div className="relative">
        <input
          id={id}
          type="number"
          min={min}
          step={step}
          value={value ?? ""}
          aria-invalid={Boolean(error)}
          aria-describedby={error ? errorID : undefined}
          className={`${fieldStyles} pr-14 ${error ? "border-danger" : ""}`}
          onChange={(event) =>
            onChange(
              event.target.value === ""
                ? undefined
                : Number(event.target.value),
            )
          }
        />
        <span className="pointer-events-none absolute right-3 bottom-3.5 font-mono text-[0.65rem] font-semibold text-ink-muted">
          {suffix}
        </span>
      </div>
      <FieldError id={errorID} message={error} />
    </div>
  );
}

function FieldError({
  id,
  message,
}: {
  id?: string;
  message?: string;
}) {
  if (!message) {
    return null;
  }

  return (
    <p id={id} className="mt-2 text-xs font-semibold text-danger">
      {message}
    </p>
  );
}

