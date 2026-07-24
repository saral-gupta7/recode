"use client";

import { mediaTools } from "@/config/tools";
import {
  type ToolFilter,
  useDashboardStore,
} from "@/stores/use-dashboard-store";

import { ToolCard } from "./tool-card";

const filters: ReadonlyArray<{ value: ToolFilter; label: string }> = [
  { value: "all", label: "All tools" },
  { value: "image", label: "Images" },
  { value: "video", label: "Video & audio" },
];

export function ToolCatalog() {
  const toolFilter = useDashboardStore((state) => state.toolFilter);
  const setToolFilter = useDashboardStore((state) => state.setToolFilter);

  const visibleTools =
    toolFilter === "all"
      ? mediaTools
      : mediaTools.filter((tool) => tool.category === toolFilter);

  return (
    <section id={`${toolFilter === "video" ? "video" : "image"}-tools`} className="scroll-mt-24 py-14 sm:py-20">
      <div className="flex flex-col justify-between gap-6 border-b border-line pb-6 md:flex-row md:items-end">
        <div>
          <p className="text-sm text-ink-muted">
            {visibleTools.length} available tools
          </p>
          <h2 className="mt-2 text-3xl leading-tight font-semibold tracking-[-0.04em] sm:text-4xl">
            Choose a tool
          </h2>
        </div>

        <div
          role="group"
          aria-label="Filter tools"
          className="flex w-full rounded-md border border-line bg-surface p-1 md:w-auto"
        >
          {filters.map((filter) => {
            const active = toolFilter === filter.value;

            return (
              <button
                key={filter.value}
                type="button"
                aria-pressed={active}
                className={`min-h-10 flex-1 px-4 font-mono text-[0.68rem] font-semibold uppercase tracking-[0.1em] transition-colors md:flex-none ${
                  active
                    ? "rounded bg-ink text-surface"
                    : "rounded bg-transparent text-ink-muted hover:bg-canvas hover:text-ink"
                }`}
                onClick={() => setToolFilter(filter.value)}
              >
                {filter.label}
              </button>
            );
          })}
        </div>
      </div>

      <p className="sr-only" aria-live="polite">
        Showing {visibleTools.length} tools
      </p>

      <div className="mt-6 grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
        {visibleTools.map((tool) => (
          <ToolCard
            key={tool.slug}
            index={mediaTools.findIndex((item) => item.slug === tool.slug)}
            tool={tool}
          />
        ))}
      </div>
    </section>
  );
}
