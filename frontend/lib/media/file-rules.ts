import type { ToolCategory } from "@/config/tools";

export const MAX_UPLOAD_BYTES = 1_073_741_824;

interface FileRule {
  mimeTypes: readonly string[];
  extensions: readonly string[];
}

export const fileRulesByCategory = {
  image: {
    mimeTypes: ["image/jpeg", "image/png", "image/webp"],
    extensions: [".jpg", ".jpeg", ".png", ".webp"],
  },

  video: {
    mimeTypes: ["video/mp4", "video/quicktime", "video/webm"],
    extensions: [".mp4", ".mov", ".webm"],
  },
} satisfies Record<ToolCategory, FileRule>;

export function getAcceptAttribute(category: ToolCategory): string {
  const rules = fileRulesByCategory[category];

  return [...rules.mimeTypes, ...rules.extensions].join(",");
}
