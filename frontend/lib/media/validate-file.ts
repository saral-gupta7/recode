import type { ToolCategory } from "@/config/tools";
import { fileRulesByCategory, MAX_UPLOAD_BYTES } from "@/lib/media/file-rules";

export type FileValidationCode =
  | "empty_file"
  | "file_too_large"
  | "unsupported_file_type";

export interface FileValidationFailure {
  valid: false;
  code: FileValidationCode;
  message: string;
}

interface FileValidationSuccess {
  valid: true;
}

export type FileValidationResult =
  | FileValidationSuccess
  | FileValidationFailure;

export function validateFile(
  file: File,
  category: ToolCategory,
  maxBytes = MAX_UPLOAD_BYTES,
): FileValidationResult {
  if (file.size === 0) {
    return {
      valid: false,
      code: "empty_file",
      message: "This file is empty. Choose a file containing media.",
    };
  }

  if (file.size > maxBytes) {
    return {
      valid: false,
      code: "file_too_large",
      message: "This file is larger than the 1 GB upload limit.",
    };
  }

  const rules = fileRulesByCategory[category];
  const normalizedMimeType = file.type.trim().toLowerCase();
  const normalizedName = file.name.trim().toLowerCase();

  const matchesMimeType = rules.mimeTypes.some(
    (mimeType) => mimeType === normalizedMimeType,
  );

  const matchesExtension = rules.extensions.some((extension) =>
    normalizedName.endsWith(extension),
  );

  if (!matchesMimeType && !matchesExtension) {
    return {
      valid: false,
      code: "unsupported_file_type",
      message:
        category === "image"
          ? "Choose a JPG, PNG, or WebP image."
          : "Choose an MP4, MOV, or WebM video.",
    };
  }

  return { valid: true };
}
