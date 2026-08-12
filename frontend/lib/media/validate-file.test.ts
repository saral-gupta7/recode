import { describe, expect, it } from "vitest";

import { validateFile } from "@/lib/media/validate-file";

function file(
  name: string,
  type: string,
  contents: BlobPart = "media",
): File {
  return new File([contents], name, { type });
}

describe("validateFile", () => {
  it("accepts an image using its MIME type", () => {
    expect(validateFile(file("photo", "image/png"), "image")).toEqual({
      valid: true,
    });
  });

  it("falls back to a supported extension when a browser omits the MIME type", () => {
    expect(validateFile(file("photo.WEBP", ""), "image")).toEqual({
      valid: true,
    });
  });

  it("rejects empty files before checking their type", () => {
    const result = validateFile(file("empty.png", "image/png", ""), "image");

    expect(result).toMatchObject({ valid: false, code: "empty_file" });
  });

  it("rejects files above the configured size boundary", () => {
    const result = validateFile(file("photo.png", "image/png"), "image", 2);

    expect(result).toMatchObject({ valid: false, code: "file_too_large" });
  });

  it("rejects a video submitted to an image tool", () => {
    const result = validateFile(file("clip.mp4", "video/mp4"), "image");

    expect(result).toMatchObject({
      valid: false,
      code: "unsupported_file_type",
    });
  });
});
