import type { MetadataRoute } from "next";

import { getSiteUrl } from "@/config/site";
import { mediaTools } from "@/config/tools";

export default function sitemap(): MetadataRoute.Sitemap {
  const siteUrl = getSiteUrl();

  return [
    {
      url: siteUrl.toString(),
      changeFrequency: "weekly",
      priority: 1,
    },
    ...mediaTools.map((tool) => ({
      url: new URL(`/tools/${tool.slug}`, siteUrl).toString(),
      changeFrequency: "monthly" as const,
      priority: 0.8,
    })),
  ];
}
