import { ArrowRight } from "lucide-react";

const stages = [
  { number: "01", label: "Input", detail: "Your file" },
  { number: "02", label: "Process", detail: "Recode engine" },
  { number: "03", label: "Output", detail: "Ready to download" },
] as const;

export function ProcessRail() {
  return (
    <ol
      aria-label="How a file moves through Recode"
      className="grid overflow-hidden rounded-lg border border-line bg-surface md:grid-cols-3"
    >
      {stages.map((stage, index) => (
        <li
          key={stage.number}
          className="relative flex items-center gap-4 border-b border-line px-5 py-4 last:border-b-0 md:border-r md:border-b-0 md:last:border-r-0"
        >
          <span className="grid size-8 shrink-0 place-items-center rounded bg-canvas font-mono text-xs text-ink-muted">
            {stage.number}
          </span>
          <span>
            <span className="block text-sm font-semibold">
              {stage.label}
            </span>
            <span className="mt-0.5 block font-mono text-[0.68rem] uppercase tracking-[0.08em] text-ink-muted">
              {stage.detail}
            </span>
          </span>
          {index < stages.length - 1 && (
            <ArrowRight
              aria-hidden="true"
              className="absolute right-[-0.7rem] z-10 hidden size-5 bg-surface text-ink-muted md:block"
            />
          )}
        </li>
      ))}
    </ol>
  );
}
