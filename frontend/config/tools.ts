import {
  AudioLines,
  Blend,
  FileVideo2,
  Film,
  ImageDown,
  Maximize2,
  RefreshCw,
  Scissors,
  VolumeX,
  type LucideIcon,
} from "lucide-react";

export type ToolCategory = "image" | "video";
export type JobOperation =
  | "video_grayscale"
  | "video_extract_audio"
  | "video_remove_audio"
  | "video_convert"
  | "video_clip"
  | "image_grayscale"
  | "image_convert"
  | "image_compress"
  | "image_resize";

export interface MediaTool {
  slug: string;
  operation: JobOperation;
  category: ToolCategory;
  title: string;
  shortTitle: string;
  description: string;
  accepts: string;
  output: string;
  icon: LucideIcon;
}

export const mediaTools: readonly MediaTool[] = [
  {
    slug: "image-grayscale",
    operation: "image_grayscale",
    category: "image",
    title: "Black & white image",
    shortTitle: "Black & white",
    description: "Remove colour while preserving the structure and contrast of an image.",
    accepts: "JPG, PNG, WebP",
    output: "PNG",
    icon: Blend,
  },
  {
    slug: "image-convert",
    operation: "image_convert",
    category: "image",
    title: "Convert image format",
    shortTitle: "Change format",
    description: "Move an image between practical web and editing formats.",
    accepts: "JPG, PNG, WebP",
    output: "JPG, PNG, WebP",
    icon: RefreshCw,
  },
  {
    slug: "image-compress",
    operation: "image_compress",
    category: "image",
    title: "Compress image",
    shortTitle: "Compress",
    description: "Reduce image size with a quality level you control.",
    accepts: "JPG, PNG, WebP",
    output: "JPG",
    icon: ImageDown,
  },
  {
    slug: "image-resize",
    operation: "image_resize",
    category: "image",
    title: "Resize image",
    shortTitle: "Resize",
    description: "Set new dimensions while keeping the image useful and sharp.",
    accepts: "JPG, PNG, WebP",
    output: "JPG",
    icon: Maximize2,
  },
  {
    slug: "video-grayscale",
    operation: "video_grayscale",
    category: "video",
    title: "Black & white video",
    shortTitle: "Black & white",
    description: "Recode a colour video into a clean monochrome version.",
    accepts: "MP4, MOV, WebM",
    output: "MP4",
    icon: FileVideo2,
  },
  {
    slug: "extract-audio",
    operation: "video_extract_audio",
    category: "video",
    title: "Extract audio",
    shortTitle: "Extract audio",
    description: "Pull the audio track out of a video as a standalone file.",
    accepts: "MP4, MOV, WebM",
    output: "MP3, WAV, M4A",
    icon: AudioLines,
  },
  {
    slug: "remove-audio",
    operation: "video_remove_audio",
    category: "video",
    title: "Remove audio",
    shortTitle: "Remove audio",
    description: "Create a silent copy of a video without its audio track.",
    accepts: "MP4, MOV, WebM",
    output: "MP4",
    icon: VolumeX,
  },
  {
    slug: "video-convert",
    operation: "video_convert",
    category: "video",
    title: "Convert video format",
    shortTitle: "Change format",
    description: "Recode a video for the browser, editing, or sharing.",
    accepts: "MP4, MOV, WebM",
    output: "MP4, MOV, WebM",
    icon: Film,
  },
  {
    slug: "video-clip",
    operation: "video_clip",
    category: "video",
    title: "Clip video",
    shortTitle: "Clip video",
    description: "Cut out the exact section you need using a start time and duration.",
    accepts: "MP4, MOV, WebM",
    output: "MP4",
    icon: Scissors,
  },
] as const;

export function getToolBySlug(slug: string): MediaTool | undefined {
  return mediaTools.find((tool) => tool.slug === slug);
}
