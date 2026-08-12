import { describe, expect, it } from "vitest";

import { defaultOptionsFor, validateOptions } from "@/features/media-jobs/options";

describe("image job options", () => {
  it("provides useful defaults", () => {
    expect(defaultOptionsFor("image_grayscale")).toEqual({});
    expect(defaultOptionsFor("image_crop")).toEqual({ x: 0, y: 0, width: 800, height: 800 });
    expect(defaultOptionsFor("image_adjust")).toEqual({ brightness: 0, contrast: 0, saturation: 100 });
    expect(defaultOptionsFor("image_padding")).toMatchObject({ padding_top: 32, background: "#ffffff" });
  });

  it("validates dimensions and crop coordinates", () => {
    expect(validateOptions("image_resize", { width: 800 })).toEqual({});
    expect(validateOptions("image_resize", { width: 0 })).toHaveProperty("width");
    expect(validateOptions("image_crop", { x: -1, y: 0, width: 100, height: 100 })).toHaveProperty("x");
  });

  it("validates transform ranges", () => {
    expect(validateOptions("image_compress", { quality: 101 })).toHaveProperty("quality");
    expect(validateOptions("image_rotate", { angle: 45 })).toHaveProperty("angle");
    expect(validateOptions("image_blur", { strength: 21 })).toHaveProperty("strength");
    expect(validateOptions("image_pixelate", { block_size: 1 })).toHaveProperty("block_size");
  });

  it("validates padding and safe colours", () => {
    expect(validateOptions("image_padding", { padding_top: 1, padding_right: 0, padding_bottom: 0, padding_left: 0, background: "#112233" })).toEqual({});
    expect(validateOptions("image_padding", { padding_top: 1, padding_right: 0, padding_bottom: 0, padding_left: 0, background: "url(x)" })).toHaveProperty("background");
  });
});
