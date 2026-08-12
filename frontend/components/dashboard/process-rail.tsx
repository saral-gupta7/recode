import { ArrowRight } from "@phosphor-icons/react/ssr";

const stages = [
  { number: "01", label: "Input", detail: "Your file" },
  { number: "02", label: "Process", detail: "Recode engine" },
  { number: "03", label: "Output", detail: "Ready to download" },
] as const;

export function ProcessRail() {
  return (
    <ol
      aria-label="How a file moves through Recode"
      className="grid overflow-hidden rounded-md border border-line bg-surface sm:grid-cols-3"
    >
      {stages.map((stage, index) => (
        <li
          key={stage.number}
          className="relative flex items-center gap-3 border-b border-line px-4 py-3 last:border-b-0 sm:border-r sm:border-b-0 sm:last:border-r-0"
        >
          <span className="grid size-7 shrink-0 place-items-center rounded bg-accent font-mono text-[0.65rem] font-semibold text-accent-strong">
            {stage.number}
          </span>
          <span>
            <span className="block text-xs font-semibold">
              {stage.label}
            </span>
            <span className="mt-0.5 block text-[0.65rem] text-ink-muted">
              {stage.detail}
            </span>
          </span>
          {index < stages.length - 1 && (
            <ArrowRight
              aria-hidden="true"
              className="absolute right-[-0.6rem] z-10 hidden size-4 bg-surface text-ink-muted sm:block"
            />
          )}
        </li>
      ))}
    </ol>
  );
}
