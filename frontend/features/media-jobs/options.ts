import type { JobOperation } from "@/config/tools";

export type FlipDirection = "horizontal" | "vertical";
export type ThumbnailPreset = "square" | "preview" | "social";
export interface JobOptions {
  format?: string;
  width?: number;
  height?: number;
  quality?: number;
  x?: number;
  y?: number;
  angle?: number;
  flip_direction?: FlipDirection;
  preset?: ThumbnailPreset;
  brightness?: number;
  contrast?: number;
  saturation?: number;
  strength?: number;
  block_size?: number;
  padding_top?: number;
  padding_right?: number;
  padding_bottom?: number;
  padding_left?: number;
  background?: string;
}

export type OptionErrors = Partial<Record<keyof JobOptions, string>>;

export function defaultOptionsFor(operation: JobOperation): JobOptions {
  switch (operation) {
    case "image_convert": return { format: "png" };
    case "image_compress": return { quality: 80 };
    case "image_resize": return { width: 1280 };
    case "image_crop": return { x: 0, y: 0, width: 800, height: 800 };
    case "image_rotate": return { angle: 90 };
    case "image_flip": return { flip_direction: "horizontal" };
    case "image_thumbnail": return { preset: "square" };
    case "image_adjust": return { brightness: 0, contrast: 0, saturation: 100 };
    case "image_blur":
    case "image_sharpen": return { strength: 2 };
    case "image_pixelate": return { block_size: 12 };
    case "image_padding": return { padding_top: 32, padding_right: 32, padding_bottom: 32, padding_left: 32, background: "#ffffff" };
    default: return {};
  }
}

export function validateOptions(operation: JobOperation, options: JobOptions): OptionErrors {
  const errors: OptionErrors = {};
  if (operation === "image_resize" && !dimension(options.width) && !dimension(options.height)) {
    errors.width = errors.height = "Enter a positive width or height (maximum 16,384).";
  }
  if (operation === "image_crop") {
    if (!nonNegativeInteger(options.x)) errors.x = "X must be zero or greater.";
    if (!nonNegativeInteger(options.y)) errors.y = "Y must be zero or greater.";
    if (!dimension(options.width)) errors.width = "Enter a valid crop width.";
    if (!dimension(options.height)) errors.height = "Enter a valid crop height.";
  }
  if (operation === "image_compress" && !betweenInteger(options.quality, 1, 100)) errors.quality = "Quality must be between 1 and 100.";
  if (operation === "image_rotate" && ![90, 180, 270].includes(options.angle ?? 0)) errors.angle = "Choose 90, 180, or 270 degrees.";
  if (operation === "image_flip" && !["horizontal", "vertical"].includes(options.flip_direction ?? "")) errors.flip_direction = "Choose a flip direction.";
  if (operation === "image_thumbnail" && !["square", "preview", "social"].includes(options.preset ?? "")) errors.preset = "Choose a thumbnail preset.";
  if (operation === "image_adjust") {
    if (!betweenInteger(options.brightness, -100, 100)) errors.brightness = "Use a value from -100 to 100.";
    if (!betweenInteger(options.contrast, -100, 100)) errors.contrast = "Use a value from -100 to 100.";
    if (!betweenInteger(options.saturation, 0, 200)) errors.saturation = "Use a value from 0 to 200.";
  }
  if (operation === "image_blur" && !between(options.strength, 0.1, 20)) errors.strength = "Strength must be between 0.1 and 20.";
  if (operation === "image_sharpen" && !between(options.strength, 0.1, 5)) errors.strength = "Strength must be between 0.1 and 5.";
  if (operation === "image_pixelate" && !betweenInteger(options.block_size, 2, 100)) errors.block_size = "Block size must be between 2 and 100.";
  if (operation === "image_padding") {
    const keys = ["padding_top", "padding_right", "padding_bottom", "padding_left"] as const;
    for (const key of keys) if (!nonNegativeInteger(options[key]) || (options[key] ?? 0) > 16384) errors[key] = "Use 0 to 16,384 pixels.";
    if (keys.every((key) => options[key] === 0)) errors.padding_top = "Add padding on at least one side.";
    if (!/^(#[0-9a-f]{6}|[a-z]{3,20})$/i.test(options.background ?? "")) errors.background = "Use a six-digit hex colour or colour name.";
  }
  return errors;
}

function dimension(value?: number) { return betweenInteger(value, 1, 16384); }
function nonNegativeInteger(value?: number) { return Number.isInteger(value) && (value ?? -1) >= 0; }
function betweenInteger(value: number | undefined, min: number, max: number) { return Number.isInteger(value) && (value ?? min - 1) >= min && (value ?? max + 1) <= max; }
function between(value: number | undefined, min: number, max: number) { return Number.isFinite(value) && (value ?? min - 1) >= min && (value ?? max + 1) <= max; }
