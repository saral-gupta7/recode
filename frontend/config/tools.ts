import {
  CircleHalf,
  Crop,
  Drop,
  Eraser,
  FlipHorizontal,
  FrameCorners,
  GridFour,
  ImageSquare,
  ImagesSquare,
  ArrowsOut,
  ArrowsClockwise,
  SlidersHorizontal,
  Sparkle,
} from "@phosphor-icons/react/ssr";
import type { Icon } from "@phosphor-icons/react";

export type ToolCategory = "image";
export type JobOperation =
  | "image_grayscale"
  | "image_convert"
  | "image_compress"
  | "image_resize"
  | "image_crop"
  | "image_rotate"
  | "image_flip"
  | "image_thumbnail"
  | "image_strip_metadata"
  | "image_adjust"
  | "image_blur"
  | "image_sharpen"
  | "image_pixelate"
  | "image_padding";

export interface MediaTool {
  slug: string;
  operation: JobOperation;
  category: ToolCategory;
  title: string;
  shortTitle: string;
  description: string;
  accepts: string;
  output: string;
  icon: Icon;
}

export const mediaTools: readonly MediaTool[] = [
  tool("image-grayscale", "image_grayscale", "Black & white image", "Black & white", "Remove colour while preserving structure and contrast.", "PNG", CircleHalf),
  tool("image-convert", "image_convert", "Convert image format", "Change format", "Convert an image to JPG, PNG, or WebP.", "JPG, PNG, WebP", ArrowsClockwise),
  tool("image-compress", "image_compress", "Compress image", "Compress", "Reduce image size with a quality level you control.", "JPG", ImageSquare),
  tool("image-resize", "image_resize", "Resize image", "Resize", "Set exact dimensions or preserve the original aspect ratio.", "JPG", ArrowsOut),
  tool("image-crop", "image_crop", "Crop image", "Crop", "Crop with coordinates measured in natural image pixels.", "PNG", Crop),
  tool("image-rotate", "image_rotate", "Rotate image", "Rotate", "Rotate clockwise by 90, 180, or 270 degrees.", "PNG", ArrowsClockwise),
  tool("image-flip", "image_flip", "Flip image", "Flip", "Mirror an image horizontally or vertically.", "PNG", FlipHorizontal),
  tool("image-thumbnail", "image_thumbnail", "Generate thumbnail", "Thumbnail", "Create a square, preview, or social sharing thumbnail.", "JPG", ImagesSquare),
  tool("image-strip-metadata", "image_strip_metadata", "Strip image metadata", "Strip metadata", "Remove EXIF, GPS, and other embedded metadata.", "JPG", Eraser),
  tool("image-adjust", "image_adjust", "Adjust image", "Adjust colour", "Tune brightness, contrast, and saturation together.", "PNG", SlidersHorizontal),
  tool("image-blur", "image_blur", "Blur image", "Blur", "Apply a controllable Gaussian blur.", "PNG", Drop),
  tool("image-sharpen", "image_sharpen", "Sharpen image", "Sharpen", "Increase edge definition with adjustable strength.", "PNG", Sparkle),
  tool("image-pixelate", "image_pixelate", "Pixelate image", "Pixelate", "Create a block-pixel treatment at the chosen size.", "PNG", GridFour),
  tool("image-padding", "image_padding", "Add canvas padding", "Add padding", "Extend the canvas with independent spacing and a background colour.", "PNG", FrameCorners),
] as const;

function tool(slug: string, operation: JobOperation, title: string, shortTitle: string, description: string, output: string, icon: Icon): MediaTool {
  return { slug, operation, category: "image", title, shortTitle, description, accepts: "JPG, PNG, WebP", output, icon };
}

export function getToolBySlug(slug: string): MediaTool | undefined {
  return mediaTools.find((item) => item.slug === slug);
}
