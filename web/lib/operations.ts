export const OPERATIONS = {
  VIDEO_TO_GIF: "VIDEO_TO_GIF",
  BLACK_AND_WHITE: "BLACK_AND_WHITE",
  EXTRACT_AUDIO: "EXTRACT_AUDIO",
  FRAME_EXTRACT: "FRAME_EXTRACT",
} as const;

export type Operation = (typeof OPERATIONS)[keyof typeof OPERATIONS];

export const OPERATION_OUTPUT: Record<Operation, { extension: string; mimeType: string; label: string; kind: "image" | "video" | "audio" }> = {
  VIDEO_TO_GIF: {
    extension: "gif",
    mimeType: "image/gif",
    label: "GIF",
    kind: "image",
  },
  BLACK_AND_WHITE: {
    extension: "mp4",
    mimeType: "video/mp4",
    label: "Video",
    kind: "video",
  },
  EXTRACT_AUDIO: {
    extension: "mp3",
    mimeType: "audio/mpeg",
    label: "Audio",
    kind: "audio",
  },
  FRAME_EXTRACT: {
    extension: "png",
    mimeType: "image/png",
    label: "Frame",
    kind: "image",
  },
};

export function isOperation(value: string): value is Operation {
  return value in OPERATIONS;
}
