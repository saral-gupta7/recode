"use client";

import { create } from "zustand";

import type { ToolCategory } from "@/config/tools";

export type ToolFilter = "all" | ToolCategory;

interface DashboardState {
  mobileNavigationOpen: boolean;
  toolFilter: ToolFilter;
  openMobileNavigation: () => void;
  closeMobileNavigation: () => void;
  toggleMobileNavigation: () => void;
  setToolFilter: (filter: ToolFilter) => void;
}

export const useDashboardStore = create<DashboardState>((set) => ({
  mobileNavigationOpen: false,
  toolFilter: "all",
  openMobileNavigation: () => set({ mobileNavigationOpen: true }),
  closeMobileNavigation: () => set({ mobileNavigationOpen: false }),
  toggleMobileNavigation: () =>
    set((state) => ({
      mobileNavigationOpen: !state.mobileNavigationOpen,
    })),
  setToolFilter: (toolFilter) => set({ toolFilter }),
}));

