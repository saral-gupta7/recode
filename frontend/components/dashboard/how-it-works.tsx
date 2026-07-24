import { Clock3, Download, Upload } from "lucide-react";

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
    <section className="overflow-hidden rounded-lg border border-line bg-surface">
      <div className="border-b border-line px-5 py-5 sm:px-8">
        <p className="text-sm font-medium text-ink-muted">
          One file. Three steps.
        </p>
      </div>
      <ol className="grid md:grid-cols-3">
        {steps.map((step) => {
          const Icon = step.icon;

          return (
            <li
              key={step.number}
              className="border-b border-line p-6 last:border-b-0 md:border-r md:border-b-0 md:last:border-r-0 sm:p-8"
            >
              <div className="flex items-center justify-between">
                <span className="font-mono text-sm text-ink-muted">
                  /{step.number}
                </span>
                <Icon aria-hidden="true" className="size-5 text-ink-muted" />
              </div>
              <h2 className="mt-10 text-xl font-semibold tracking-[-0.025em]">
                {step.title}
              </h2>
              <p className="mt-3 text-sm leading-6 text-ink-muted">
                {step.description}
              </p>
            </li>
          );
        })}
      </ol>
    </section>
  );
}
