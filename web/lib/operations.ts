export const OPERATIONS = {
  BLACK_AND_WHITE: "BLACK_AND_WHITE",
  REMOVE_AUDIO: "REMOVE_AUDIO",
  EXTRACT_AUDIO: "EXTRACT_AUDIO",
  FRAME_EXTRACT: "FRAME_EXTRACT",
  TRIM_30S: "TRIM_30S",
  COMPRESS_VIDEO: "COMPRESS_VIDEO",
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
  TRIM_30S: {
    extension: "mp4",
    mimeType: "video/mp4",
    label: "Video",
    kind: "video",
  },
  COMPRESS_VIDEO: {
    extension: "mp4",
    mimeType: "video/mp4",
    label: "Video",
    kind: "video",
  },
};

export function isOperation(value: string): value is Operation {
  return value in OPERATIONS;
}
