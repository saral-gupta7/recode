import { LoaderCircle } from "lucide-react";

export default function Loading() {
  return (
    <main
      id="main-content"
      className="grid min-h-[70svh] place-items-center px-6"
    >
      <div className="text-center">
        <LoaderCircle
          aria-hidden="true"
          className="mx-auto size-10 animate-spin text-accent-strong"
        />
        <p className="mt-5 font-mono text-xs font-semibold uppercase tracking-[0.14em] text-ink-muted">
          Preparing workspace
        </p>
      </div>
    </main>
  );
}
