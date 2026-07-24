"use client";

import { useEffect } from "react";
import Link from "next/link";
import { ArrowUpRight, X } from "lucide-react";

import { BrandMark } from "@/components/brand/brand-mark";
import { ThemeToggle } from "@/components/ui/theme-toggle";
import { mediaTools, type ToolCategory } from "@/config/tools";
import { useDashboardStore } from "@/stores/use-dashboard-store";

const groups: ReadonlyArray<{
  category: ToolCategory;
  label: string;
  number: string;
}> = [
  { category: "image", label: "Image tools", number: "01" },
  { category: "video", label: "Video & audio", number: "02" },
];

export function Sidebar() {
  const mobileNavigationOpen = useDashboardStore(
    (state) => state.mobileNavigationOpen,
  );
  const closeMobileNavigation = useDashboardStore(
    (state) => state.closeMobileNavigation,
  );
  const setToolFilter = useDashboardStore((state) => state.setToolFilter);

  useEffect(() => {
    if (!mobileNavigationOpen) {
      return;
    }

    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = "hidden";

    return () => {
      document.body.style.overflow = previousOverflow;
    };
  }, [mobileNavigationOpen]);

  function showCategory(category: ToolCategory) {
    setToolFilter(category);
    closeMobileNavigation();
  }

  return (
    <>
      <button
        type="button"
        aria-label="Close navigation"
        className={`fixed inset-0 z-40 bg-black/55 backdrop-blur-[2px] transition-opacity duration-200 lg:hidden ${
          mobileNavigationOpen
            ? "pointer-events-auto opacity-100"
            : "pointer-events-none opacity-0"
        }`}
        onClick={closeMobileNavigation}
      />

      <aside
        id="primary-navigation"
        aria-label="Tool navigation"
        className={`fixed inset-y-0 left-0 z-50 flex w-[19rem] flex-col border-r border-line bg-surface transition-transform duration-300 ease-workshop lg:translate-x-0 ${
          mobileNavigationOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <div className="flex h-20 shrink-0 items-center justify-between border-b border-line px-6">
          <BrandMark />
          <button
            type="button"
            aria-label="Close navigation"
            className="grid size-11 place-items-center rounded-md border-2 border-transparent hover:border-ink hover:bg-accent lg:hidden"
            onClick={closeMobileNavigation}
          >
            <X aria-hidden="true" className="size-5" />
          </button>
        </div>

        <nav className="flex-1 overflow-y-auto px-4 py-6">
          <Link
            href="/"
            className="mb-7 flex items-center justify-between rounded-md bg-accent px-4 py-3 text-sm font-semibold"
            onClick={() => {
              setToolFilter("all");
              closeMobileNavigation();
            }}
          >
            All tools
            <ArrowUpRight aria-hidden="true" className="size-4" />
          </Link>

          <div className="space-y-8">
            {groups.map((group) => (
              <section key={group.category}>
                <Link
                  href={`/#${group.category}-tools`}
                  className="mb-3 flex items-center gap-3 px-2 font-mono text-[0.68rem] font-semibold uppercase tracking-[0.16em] text-ink-muted hover:text-ink"
                  onClick={() => showCategory(group.category)}
                >
                  <span>{group.number}</span>
                  <span className="h-px flex-1 bg-line-strong" />
                  <span>{group.label}</span>
                </Link>

                <ul className="space-y-1">
                  {mediaTools
                    .filter((tool) => tool.category === group.category)
                    .map((tool) => {
                      const Icon = tool.icon;

                      return (
                        <li key={tool.slug}>
                          <Link
                            href={`/tools/${tool.slug}`}
                            className="group flex min-h-11 items-center gap-3 rounded-md px-3 py-2 text-sm font-medium hover:bg-canvas"
                            onClick={closeMobileNavigation}
                          >
                            <Icon
                              aria-hidden="true"
                              className="size-4 shrink-0"
                            />
                            <span>{tool.shortTitle}</span>
                          </Link>
                        </li>
                      );
                    })}
                </ul>
              </section>
            ))}
          </div>
        </nav>

        <div className="border-t border-line p-4">
          <div className="mb-3 flex items-center justify-between">
            <span className="font-mono text-[0.65rem] font-semibold uppercase tracking-[0.12em] text-ink-muted">
              Appearance
            </span>
            <ThemeToggle />
          </div>
        </div>
      </aside>
    </>
  );
}
