import {
  Clock as Clock3,
  DownloadSimple as Download,
  UploadSimple as Upload,
} from "@phosphor-icons/react/ssr";

const steps = [
  {
    number: "01",
    title: "Choose and upload",
    description: "Pick one focused tool and add the file you want to change.",
    icon: Upload,
  },
  {
    number: "02",
    title: "Let Recode work",
    description: "The job runs away from the request, so the interface stays responsive.",
    icon: Clock3,
  },
  {
    number: "03",
    title: "Preview and download",
    description: "Check the result in the browser, then save it to your device.",
    icon: Download,
  },
] as const;

export function HowItWorks() {
  return (
    <section className="mb-6 overflow-hidden rounded-md border border-line bg-surface">
      <div className="border-b border-line px-5 py-3.5">
        <p className="text-xs font-semibold text-ink-muted">
          One file. Three steps.
        </p>
      </div>
      <ol className="grid md:grid-cols-3">
        {steps.map((step) => {
          const Icon = step.icon;

          return (
            <li
              key={step.number}
              className="border-b border-line p-5 last:border-b-0 md:border-r md:border-b-0 md:last:border-r-0"
            >
              <div className="flex items-center justify-between">
                <span className="font-mono text-[0.65rem] font-semibold text-ink-muted">
                  /{step.number}
                </span>
                <Icon aria-hidden="true" className="size-4 text-accent-strong" />
              </div>
              <h2 className="mt-5 text-sm font-semibold tracking-[-0.015em]">
                {step.title}
              </h2>
              <p className="mt-2 text-xs leading-5 text-ink-muted">
                {step.description}
              </p>
            </li>
          );
        })}
      </ol>
    </section>
  );
}
