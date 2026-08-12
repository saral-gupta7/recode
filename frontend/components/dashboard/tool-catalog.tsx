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
];

export function ToolCatalog() {
  const toolFilter = useDashboardStore((state) => state.toolFilter);
  const setToolFilter = useDashboardStore((state) => state.setToolFilter);

  const visibleTools =
    toolFilter === "all"
      ? mediaTools
      : mediaTools.filter((tool) => tool.category === toolFilter);

  return (
    <section id="image-tools" className="scroll-mt-20 py-10 sm:py-12">
      <div className="flex flex-col justify-between gap-4 border-b border-line pb-4 sm:flex-row sm:items-end">
        <div>
          <p className="text-xs font-medium text-accent-strong">
            {visibleTools.length} available tools
          </p>
          <h2 className="mt-1 text-2xl leading-tight font-semibold tracking-[-0.04em] sm:text-3xl">
            What do you need?
          </h2>
        </div>

        <div
          role="group"
          aria-label="Filter tools"
          className="flex w-full rounded-md border border-line bg-surface p-1 sm:w-auto"
        >
          {filters.map((filter) => {
            const active = toolFilter === filter.value;

            return (
              <button
                key={filter.value}
                type="button"
                aria-pressed={active}
                className={`min-h-8 flex-1 px-3 text-xs font-semibold sm:flex-none ${
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

      <div className="mt-4 grid gap-2.5 sm:grid-cols-2">
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
