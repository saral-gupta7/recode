"use client";

import { List as Menu } from "@phosphor-icons/react/ssr";

import { BrandMark } from "@/components/brand/brand-mark";
import { ThemeToggle } from "@/components/ui/theme-toggle";
import { useDashboardStore } from "@/stores/use-dashboard-store";

export function MobileHeader() {
  const openMobileNavigation = useDashboardStore(
    (state) => state.openMobileNavigation,
  );

  return (
    <header className="sticky top-0 z-30 flex h-14 items-center justify-between border-b border-line bg-surface/90 px-4 backdrop-blur-xl lg:hidden">
      <BrandMark />
      <div className="flex items-center gap-2">
        <ThemeToggle />
        <button
          type="button"
          aria-controls="primary-navigation"
          aria-label="Open tool navigation"
          className="grid size-9 place-items-center rounded-md border border-line bg-canvas hover:border-line-strong hover:bg-accent"
          onClick={openMobileNavigation}
        >
          <Menu aria-hidden="true" className="size-4" />
        </button>
      </div>
    </header>
  );
}
