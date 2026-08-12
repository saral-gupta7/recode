"use client";

import type { JobOperation } from "@/config/tools";
import type { JobOptions, OptionErrors } from "@/features/media-jobs/options";

interface Props { operation: JobOperation; options: JobOptions; errors: OptionErrors; onChange: (options: JobOptions) => void }
const fieldStyles = "mt-1.5 min-h-10 w-full rounded-md border border-line-strong bg-surface px-3 text-xs outline-none hover:border-line-strong focus:border-focus";
const labelStyles = "block text-[0.68rem] font-semibold text-ink-muted";

export function ToolOptions({ operation, options, errors, onChange }: Props) {
  const set = <K extends keyof JobOptions>(key: K, value: JobOptions[K]) => onChange({ ...options, [key]: value });

  if (operation === "image_convert") return <SelectField label="Output format" value={options.format ?? "png"} values={["png", "jpg", "webp"]} onChange={(value) => set("format", value)} />;
  if (operation === "image_compress") return <SliderField id="quality" label="Output quality" value={options.quality ?? 80} min={1} max={100} suffix="%" error={errors.quality} onChange={(value) => set("quality", value)} />;
  if (operation === "image_resize") return <NumberGrid fields={[
    { name: "width", label: "Width", value: options.width, error: errors.width },
    { name: "height", label: "Height", value: options.height, error: errors.height },
  ]} set={set} help="Leave one dimension empty to preserve the original aspect ratio." />;
  if (operation === "image_crop") return <NumberGrid fields={[
    { name: "x", label: "X position", value: options.x, error: errors.x, min: 0 },
    { name: "y", label: "Y position", value: options.y, error: errors.y, min: 0 },
    { name: "width", label: "Crop width", value: options.width, error: errors.width },
    { name: "height", label: "Crop height", value: options.height, error: errors.height },
  ]} set={set} help="Coordinates use the image's natural pixels and must stay within its bounds." />;
  if (operation === "image_rotate") return <SelectField label="Clockwise rotation" value={String(options.angle ?? 90)} values={["90", "180", "270"]} suffix="°" onChange={(value) => set("angle", Number(value))} />;
  if (operation === "image_flip") return <SelectField label="Flip direction" value={options.flip_direction ?? "horizontal"} values={["horizontal", "vertical"]} onChange={(value) => set("flip_direction", value as "horizontal" | "vertical")} />;
  if (operation === "image_thumbnail") return <SelectField label="Thumbnail preset" value={options.preset ?? "square"} values={["square", "preview", "social"]} labels={["Square — 256 × 256", "Preview — 640 × 360", "Social — 1200 × 630"]} onChange={(value) => set("preset", value as "square" | "preview" | "social")} />;
  if (operation === "image_adjust") return <div className="space-y-4">
    <SliderField id="brightness" label="Brightness" value={options.brightness ?? 0} min={-100} max={100} error={errors.brightness} onChange={(value) => set("brightness", value)} />
    <SliderField id="contrast" label="Contrast" value={options.contrast ?? 0} min={-100} max={100} error={errors.contrast} onChange={(value) => set("contrast", value)} />
    <SliderField id="saturation" label="Saturation" value={options.saturation ?? 100} min={0} max={200} suffix="%" error={errors.saturation} onChange={(value) => set("saturation", value)} />
  </div>;
  if (operation === "image_blur" || operation === "image_sharpen") return <SliderField id="strength" label={`${operation === "image_blur" ? "Blur" : "Sharpen"} strength`} value={options.strength ?? 2} min={0.1} max={operation === "image_blur" ? 20 : 5} step={0.1} error={errors.strength} onChange={(value) => set("strength", value)} />;
  if (operation === "image_pixelate") return <SliderField id="block-size" label="Pixel block size" value={options.block_size ?? 12} min={2} max={100} suffix=" px" error={errors.block_size} onChange={(value) => set("block_size", value)} />;
  if (operation === "image_padding") return <div className="space-y-4">
    <NumberGrid fields={[
      { name: "padding_top", label: "Top", value: options.padding_top, error: errors.padding_top, min: 0 },
      { name: "padding_right", label: "Right", value: options.padding_right, error: errors.padding_right, min: 0 },
      { name: "padding_bottom", label: "Bottom", value: options.padding_bottom, error: errors.padding_bottom, min: 0 },
      { name: "padding_left", label: "Left", value: options.padding_left, error: errors.padding_left, min: 0 },
    ]} set={set} />
    <label className={labelStyles}>Background colour<input type="text" value={options.background ?? "#ffffff"} className={fieldStyles} placeholder="#ffffff" onChange={(event) => set("background", event.target.value)} /><FieldError message={errors.background} /></label>
  </div>;

  return <p className="text-sm leading-6 text-ink-muted">No additional settings are needed. Recode will use a balanced image profile.</p>;
}

interface NumberConfig { name: keyof JobOptions; label: string; value?: number; error?: string; min?: number }
function NumberGrid({ fields, set, help }: { fields: NumberConfig[]; set: <K extends keyof JobOptions>(key: K, value: JobOptions[K]) => void; help?: string }) {
  return <div className="grid gap-3 sm:grid-cols-2">{fields.map((field) => <NumberField key={field.name} {...field} onChange={(value) => set(field.name, value)} />)}{help && <p className="text-[0.68rem] leading-4 text-ink-muted sm:col-span-2">{help}</p>}</div>;
}

function SelectField({ label, value, values, labels, suffix = "", onChange }: { label: string; value: string; values: readonly string[]; labels?: readonly string[]; suffix?: string; onChange: (value: string) => void }) {
  return <label className={labelStyles}>{label}<select value={value} className={fieldStyles} onChange={(event) => onChange(event.target.value)}>{values.map((item, index) => <option key={item} value={item}>{labels?.[index] ?? `${item.toUpperCase()}${suffix}`}</option>)}</select></label>;
}

function SliderField({ id, label, value, min, max, step = 1, suffix = "", error, onChange }: { id: string; label: string; value: number; min: number; max: number; step?: number; suffix?: string; error?: string; onChange: (value: number) => void }) {
  return <div><div className="flex items-center justify-between"><label htmlFor={id} className={labelStyles}>{label}</label><output htmlFor={id} className="font-mono text-xs font-semibold text-accent-strong">{value}{suffix}</output></div><input id={id} type="range" min={min} max={max} step={step} value={value} className="mt-3 w-full accent-[var(--color-accent-strong)]" onChange={(event) => onChange(Number(event.target.value))} /><FieldError message={error} /></div>;
}

function NumberField({ name, label, value, error, min = 1, onChange }: NumberConfig & { onChange: (value: number | undefined) => void }) {
  const id = String(name);
  return <div><label htmlFor={id} className={labelStyles}>{label}</label><div className="relative"><input id={id} type="number" min={min} max={16384} step={1} value={value ?? ""} className={`${fieldStyles} pr-11 ${error ? "border-danger" : ""}`} onChange={(event) => onChange(event.target.value === "" ? undefined : Number(event.target.value))} /><span className="pointer-events-none absolute right-3 bottom-3 font-mono text-[0.58rem] font-semibold text-ink-muted">PX</span></div><FieldError message={error} /></div>;
}

function FieldError({ message }: { message?: string }) { return message ? <p className="mt-1.5 text-[0.68rem] font-semibold text-danger">{message}</p> : null; }
