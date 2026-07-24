import Link from "next/link";
import { ArrowLeft } from "lucide-react";

export default function NotFound() {
  return (
    <main
      id="main-content"
      className="grid min-h-[70svh] place-items-center px-6 py-16"
    >
      <section className="w-full max-w-3xl border-2 border-ink bg-surface p-8 sm:p-12">
        <p className="font-mono text-xs font-semibold uppercase tracking-[0.14em] text-ink-muted">
          Error / 404
        </p>
        <h1 className="mt-6 text-[clamp(3rem,10vw,3rem)] leading-[0.85] font-black tracking-[-0.07em] uppercase">
          Tool not found
        </h1>
        <p className="mt-7 max-w-lg text-lg leading-8 text-ink-muted">
          This workshop route does not exist or has moved back to the catalog.
        </p>
        <Link
          href="/"
          className="mt-8 inline-flex min-h-12 items-center gap-2 border-2 border-ink bg-accent px-5 py-3 text-sm font-bold text-[#171a17] shadow-hard"
        >
          <ArrowLeft aria-hidden="true" className="size-4" />
          Browse all tools
        </Link>
      </section>
    </main>
  );
}
