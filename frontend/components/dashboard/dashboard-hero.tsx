import { ProcessRail } from "@/components/dashboard/process-rail";

const facts = [
  { value: "14", label: "Focused tools" },
  { value: "0", label: "Accounts needed" },
  { value: "2h", label: "Private retention" },
] as const;

export function DashboardHero() {
  return (
    <section className="workshop-enter pt-8 sm:pt-10">
      <div className="grid items-end gap-6 rounded-lg border border-line bg-surface p-6 shadow-hard sm:p-8 lg:grid-cols-[minmax(0,1fr)_18rem]">
        <div>
          <p className="mb-3 text-sm font-semibold text-accent-strong">
            Private image workspace
          </p>

          <h1 className="max-w-3xl text-balance text-[clamp(2.35rem,5vw,4.6rem)] leading-[0.98] font-semibold tracking-[-0.055em]">
            Image tools, without the clutter.
          </h1>

          <p className="mt-4 max-w-xl text-pretty text-sm leading-6 text-ink-muted sm:text-base">
            Crop, resize, adjust, and export in a few clicks. No account,
            advertising, or software to install.
          </p>
        </div>
        <div className="grid grid-cols-3 gap-2 lg:grid-cols-1">
          {facts.map((fact) => (
            <div key={fact.label} className="rounded-md bg-canvas px-3 py-2.5 lg:flex lg:items-baseline lg:justify-between">
              <span className="block font-mono text-sm font-semibold text-ink">{fact.value}</span>
              <span className="mt-1 block text-[0.68rem] leading-tight text-ink-muted lg:mt-0">{fact.label}</span>
            </div>
          ))}
        </div>
      </div>

      <div className="mt-3">
        <ProcessRail />
      </div>
    </section>
  );
}
