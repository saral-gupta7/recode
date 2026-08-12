import Link from "next/link";
import { ArrowLeft } from "@phosphor-icons/react/ssr";

export default function NotFound() {
  return (
    <main
      id="main-content"
      className="grid min-h-[70svh] place-items-center px-6 py-16"
    >
      <section className="w-full max-w-xl rounded-lg border border-line bg-surface p-8 shadow-hard sm:p-10">
        <p className="text-xs font-semibold text-accent-strong">
          Error 404
        </p>
        <h1 className="mt-3 text-4xl leading-tight font-semibold tracking-[-0.045em]">
          Tool not found
        </h1>
        <p className="mt-4 max-w-lg text-sm leading-6 text-ink-muted">
          This tool does not exist or has moved back to the catalog.
        </p>
        <Link
          href="/"
          className="mt-7 inline-flex min-h-10 items-center gap-2 rounded-md bg-ink px-4 py-2 text-sm font-semibold text-canvas"
        >
          <ArrowLeft aria-hidden="true" className="size-4" />
          Browse all tools
        </Link>
      </section>
    </main>
  );
}
