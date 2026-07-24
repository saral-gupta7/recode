import { ProcessRail } from "@/components/dashboard/process-rail";

const facts = [
  { value: "09", label: "Focused tools" },
  { value: "00", label: "Accounts needed" },
  { value: "2H", label: "Result retention" },
] as const;

export function DashboardHero() {
  return (
    <section className="workshop-enter">
      <div className="grid gap-10 border-b border-line py-14 lg:grid-cols-[minmax(0,1fr)_16rem] lg:py-20">
        <div>
          <p className="mb-5 text-sm font-medium text-ink-muted">
            Private media utilities
          </p>

          <h1 className="max-w-3xl text-[clamp(2.8rem,3vw,5.8rem)] leading-[0.96] font-semibold tracking-[-0.06em]">
            Simple tools for images and video.
          </h1>

          <p className="mt-7 max-w-2xl text-base leading-7 text-ink-muted sm:text-md">
            Convert, resize, compress, clip, or extract what you need. No
            account, advertising, or software installation.
          </p>
        </div>
      </div>

      <div className="grid border-b border-line sm:grid-cols-3">
        {facts.map((fact) => (
          <div
            key={fact.label}
            className="flex items-baseline gap-3 border-b border-line py-4 last:border-b-0 sm:border-r sm:border-b-0 sm:px-5 sm:first:pl-0 sm:last:border-r-0"
          >
            <span className="font-mono text-lg font-semibold">
              {fact.value}
            </span>
            <span className="text-xs text-ink-muted">{fact.label}</span>
          </div>
        ))}
      </div>

      <div className="mt-6">
        <ProcessRail />
      </div>
    </section>
  );
}
