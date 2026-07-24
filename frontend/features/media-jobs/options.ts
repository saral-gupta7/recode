import type { JobOperation } from "@/config/tools";

export interface JobOptions {
  format?: string;
  width?: number;
  height?: number;
  quality?: number;
  start_seconds?: number;
  duration_seconds?: number;
}

export type OptionErrors = Partial<Record<keyof JobOptions, string>>;

export function defaultOptionsFor(operation: JobOperation): JobOptions {
  switch (operation) {
    case "image_convert":
      return { format: "png" };
    case "image_compress":
      return { quality: 80 };
    case "image_resize":
      return { width: 1280 };
    case "video_extract_audio":
      return { format: "mp3" };
    case "video_convert":
      return { format: "mp4" };
    case "video_clip":
      return { start_seconds: 0, duration_seconds: 10 };
    default:
      return {};
  }
}

export function validateOptions(
  operation: JobOperation,
  options: JobOptions,
): OptionErrors {
  const errors: OptionErrors = {};

  if (
    operation === "image_resize" &&
    (!positiveInteger(options.width) && !positiveInteger(options.height))
  ) {
    errors.width = "Enter a positive width or height.";
    errors.height = "Enter a positive width or height.";
  }

  if (
    operation === "image_resize" &&
    ((options.width ?? 0) > 16384 || (options.height ?? 0) > 16384)
  ) {
    errors.width = "Dimensions cannot exceed 16,384 pixels.";
  }

  if (
    operation === "image_compress" &&
    (!Number.isInteger(options.quality) ||
      (options.quality ?? 0) < 1 ||
      (options.quality ?? 0) > 100)
  ) {
    errors.quality = "Quality must be between 1 and 100.";
  }

  if (operation === "video_clip") {
    if (
      !Number.isFinite(options.start_seconds) ||
      (options.start_seconds ?? -1) < 0
    ) {
      errors.start_seconds = "Start time cannot be negative.";
    }
    if (
      !Number.isFinite(options.duration_seconds) ||
      (options.duration_seconds ?? 0) <= 0 ||
      (options.duration_seconds ?? 0) > 21600
    ) {
      errors.duration_seconds =
        "Duration must be greater than zero and no more than six hours.";
    }
  }

  return errors;
}

function positiveInteger(value: number | undefined): boolean {
  return (
    value !== undefined &&
    Number.isInteger(value) &&
    value > 0 &&
    value <= 16384
  );
}

