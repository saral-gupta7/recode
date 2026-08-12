"use client";

import { useEffect } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { SquaresFour as LayoutGrid, X } from "@phosphor-icons/react/ssr";

import { BrandMark } from "@/components/brand/brand-mark";
import { ThemeToggle } from "@/components/ui/theme-toggle";
import { mediaTools, type ToolCategory } from "@/config/tools";
import { useDashboardStore } from "@/stores/use-dashboard-store";

const groups: ReadonlyArray<{
  category: ToolCategory;
  label: string;
}> = [
  { category: "image", label: "Image tools" },
];

export function Sidebar() {
  const pathname = usePathname();
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
        className={`fixed inset-y-0 left-0 z-50 flex w-64 flex-col border-r border-line bg-surface/95 shadow-soft backdrop-blur-xl transition-transform duration-300 ease-workshop lg:translate-x-0 lg:shadow-none ${
          mobileNavigationOpen ? "translate-x-0" : "-translate-x-full"
        }`}
      >
        <div className="flex h-16 shrink-0 items-center justify-between border-b border-line px-5">
          <BrandMark />
          <button
            type="button"
            aria-label="Close navigation"
            className="grid size-9 place-items-center rounded-md hover:bg-accent lg:hidden"
            onClick={closeMobileNavigation}
          >
            <X aria-hidden="true" className="size-5" />
          </button>
        </div>

        <nav className="flex-1 overflow-y-auto px-3 py-4">
          <Link
            href="/"
            aria-current={pathname === "/" ? "page" : undefined}
            className={`mb-4 flex min-h-10 items-center gap-3 rounded-md px-3 text-sm font-semibold ${pathname === "/" ? "bg-accent text-accent-strong" : "text-ink-muted hover:bg-canvas hover:text-ink"}`}
            onClick={() => {
              setToolFilter("all");
              closeMobileNavigation();
            }}
          >
            <LayoutGrid aria-hidden="true" className="size-4" />
            All tools
          </Link>

          <div className="space-y-5">
            {groups.map((group) => (
              <section key={group.category}>
                <Link
                  href={`/#${group.category}-tools`}
                  className="mb-2 flex items-center gap-3 px-3 text-xs font-semibold text-ink-muted hover:text-ink"
                  onClick={() => showCategory(group.category)}
                >
                  <span>{group.label}</span>
                  <span className="ml-auto rounded bg-canvas px-1.5 py-0.5 font-mono text-[0.6rem]">{mediaTools.length}</span>
                </Link>

                <ul className="space-y-0.5">
                  {mediaTools
                    .filter((tool) => tool.category === group.category)
                    .map((tool) => {
                      const Icon = tool.icon;
                      const href = `/tools/${tool.slug}`;
                      const active = pathname === href;

                      return (
                        <li key={tool.slug}>
                          <Link
                            href={href}
                            aria-current={active ? "page" : undefined}
                            className={`group flex min-h-9 items-center gap-2.5 rounded-md px-3 py-1.5 text-[0.82rem] font-medium ${active ? "bg-accent text-accent-strong" : "text-ink-muted hover:bg-canvas hover:text-ink"}`}
                            onClick={closeMobileNavigation}
                          >
                            <Icon
                              aria-hidden="true"
                              weight="bold"
                              className="size-3.5 shrink-0"
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

        <div className="border-t border-line p-3">
          <div className="flex items-center justify-between rounded-md bg-canvas px-3 py-2">
            <span className="text-xs font-medium text-ink-muted">Appearance</span>
            <ThemeToggle />
          </div>
        </div>
      </aside>
    </>
  );
}
