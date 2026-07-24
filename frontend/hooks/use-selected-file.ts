"use client";

import { useCallback, useEffect, useRef, useState } from "react";

import type { ToolCategory } from "@/config/tools";
import {
  type FileValidationFailure,
  validateFile,
} from "@/lib/media/validate-file";

export interface SelectedFile {
  file: File;
  previewUrl: string;
}

export function useSelectedFile(category: ToolCategory) {
  const [selectedFile, setSelectedFile] = useState<SelectedFile | null>(null);

  const [error, setError] = useState<FileValidationFailure | null>(null);

  const previewUrlRef = useRef<string | null>(null);

  const releasePreviewUrl = useCallback(() => {
    if (!previewUrlRef.current) {
      return;
    }

    URL.revokeObjectURL(previewUrlRef.current);
    previewUrlRef.current = null;
  }, []);

  const selectFile = useCallback(
    (candidate: File | null): boolean => {
      if (!candidate) {
        return false;
      }

      const validation = validateFile(candidate, category);

      if (!validation.valid) {
        setError(validation);
        return false;
      }

      releasePreviewUrl();

      const previewUrl = URL.createObjectURL(candidate);
      previewUrlRef.current = previewUrl;

      setSelectedFile({
        file: candidate,
        previewUrl,
      });

      setError(null);
      return true;
    },
    [category, releasePreviewUrl],
  );

  const clearFile = useCallback(() => {
    releasePreviewUrl();
    setSelectedFile(null);
    setError(null);
  }, [releasePreviewUrl]);

  useEffect(() => {
    return releasePreviewUrl;
  }, [releasePreviewUrl]);

  return {
    selectedFile,
    error,
    selectFile,
    clearFile,
  };
}
