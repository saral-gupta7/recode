"use client";

import { useEffect } from "react";

import { useWorkspaceStore } from "@/stores/use-workspace-store";

export function WorkspaceProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  useEffect(() => {
    const markHydrated = () => {
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
