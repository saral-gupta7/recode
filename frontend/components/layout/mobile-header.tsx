"use client";

import { Menu } from "lucide-react";

import { BrandMark } from "@/components/brand/brand-mark";
import { ThemeToggle } from "@/components/ui/theme-toggle";
import { useDashboardStore } from "@/stores/use-dashboard-store";

export function MobileHeader() {
  const openMobileNavigation = useDashboardStore(
    (state) => state.openMobileNavigation,
  );

  return (
    <header className="sticky top-0 z-30 flex h-16 items-center justify-between border-b border-line bg-surface/95 px-4 backdrop-blur lg:hidden">
      <BrandMark />
      <div className="flex items-center gap-2">
        <ThemeToggle />
        <button
          type="button"
          aria-controls="primary-navigation"
          aria-label="Open tool navigation"
          className="grid size-11 place-items-center rounded-md border border-line bg-surface hover:bg-canvas"
          onClick={openMobileNavigation}
        >
          <Menu aria-hidden="true" className="size-5" />
        </button>
      </div>
    </header>
  );
}
