import { describe, expect, it } from "vitest";

import {
  defaultOptionsFor,
  validateOptions,
} from "@/features/media-jobs/options";

describe("defaultOptionsFor", () => {
  it("provides useful defaults only when an operation needs options", () => {
    expect(defaultOptionsFor("image_grayscale")).toEqual({});
    expect(defaultOptionsFor("image_compress")).toEqual({ quality: 80 });
    expect(defaultOptionsFor("video_clip")).toEqual({
      start_seconds: 0,
      duration_seconds: 10,
    });
  });
});

describe("validateOptions", () => {
  it("allows resize with either one or both dimensions", () => {
    expect(validateOptions("image_resize", { width: 800 })).toEqual({});
    expect(validateOptions("image_resize", { height: 600 })).toEqual({});
  });

  it("requires at least one valid resize dimension", () => {
    expect(
      validateOptions("image_resize", { width: 0, height: -1 }),
    ).toMatchObject({
      width: expect.any(String),
      height: expect.any(String),
    });
  });

  it("rejects invalid quality boundaries", () => {
    expect(validateOptions("image_compress", { quality: 0 })).toHaveProperty(
      "quality",
    );
    expect(validateOptions("image_compress", { quality: 101 })).toHaveProperty(
      "quality",
    );
    expect(validateOptions("image_compress", { quality: 75 })).toEqual({});
  });

  it("rejects negative clip starts and non-positive durations", () => {
    expect(
      validateOptions("video_clip", {
        start_seconds: -1,
        duration_seconds: 0,
      }),
    ).toMatchObject({
      start_seconds: expect.any(String),
      duration_seconds: expect.any(String),
    });
  });
});
