"use client";

import { useEffect } from "react";

import { getToolBySlug } from "@/config/tools";
import { useWorkspaceStore } from "@/stores/use-workspace-store";

export function WorkspaceProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  useEffect(() => {
    const markHydrated = () => {
      const state = useWorkspaceStore.getState();
      if (state.activeJob && !getToolBySlug(state.activeJob.toolSlug)) {
        state.clearActiveJob();
      }
      useWorkspaceStore.setState({ hydrated: true });
    };

    const unsubscribe =
      useWorkspaceStore.persist.onFinishHydration(markHydrated);

    void Promise.resolve(useWorkspaceStore.persist.rehydrate())
      .then(markHydrated)
      .catch(markHydrated);

    return unsubscribe;
  }, []);

  return children;
}
