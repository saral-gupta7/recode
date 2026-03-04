export const OPERATIONS = {
  BLACK_AND_WHITE: "BLACK_AND_WHITE",
  REMOVE_AUDIO: "REMOVE_AUDIO",
  FRAME_EXTRACT: "FRAME_EXTRACT",
} as const;

export type Operation = (typeof OPERATIONS)[keyof typeof OPERATIONS];

export const OPERATION_OUTPUT: Record<
  Operation,
  { extension: string; mimeType: string; label: string; kind: "image" | "video" | "audio" }
> = {
  BLACK_AND_WHITE: {
    extension: "mp4",
    mimeType: "video/mp4",
    label: "Video",
    kind: "video",
  },
  REMOVE_AUDIO: {
    extension: "mp4",
    mimeType: "video/mp4",
    label: "Video",
    kind: "video",
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
