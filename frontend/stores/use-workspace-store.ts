"use client";

import { create } from "zustand";
import {
  createJSONStorage,
  persist,
} from "zustand/middleware";

import type { JobOperation } from "@/config/tools";

export interface ActiveJobCredentials {
  id: string;
  ownerToken: string;
  toolSlug: string;
  operation: JobOperation;
  createdAt: string;
}

interface WorkspaceState {
  activeJob: ActiveJobCredentials | null;
  hydrated: boolean;
  setActiveJob: (job: ActiveJobCredentials) => void;
  clearActiveJob: () => void;
  setHydrated: (hydrated: boolean) => void;
}

export const useWorkspaceStore = create<WorkspaceState>()(
  persist(
    (set) => ({
      activeJob: null,
      hydrated: false,
      setActiveJob: (activeJob) => set({ activeJob }),
      clearActiveJob: () => set({ activeJob: null }),
      setHydrated: (hydrated) => set({ hydrated }),
    }),
    {
      name: "recode-active-job",
      storage: createJSONStorage(() => sessionStorage),
      partialize: (state) => ({ activeJob: state.activeJob }),
      skipHydration: true,
    },
  ),
);
