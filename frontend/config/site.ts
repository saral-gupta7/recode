export const siteConfig = {
  name: "Recode",
  title: "Recode — Private, account-free image tools",
  description:
    "Convert, compress, resize, crop, adjust, and export images with 14 focused tools. No account, advertising, or software installation required.",
  shortDescription:
    "Private, account-free tools for converting and processing images.",
  keywords: [
    "image tools",
    "image converter",
    "image compressor",
    "image resizer",
    "image cropper",
    "image editor",
    "convert image online",
    "compress image online",
    "resize image online",
    "private image processing",
  ],
} as const;

export function getSiteUrl(): URL {
  const configuredUrl = process.env.RECODE_SITE_URL?.trim();

  if (!configuredUrl) {
    return new URL("http://localhost:3000");
  }

  return new URL(configuredUrl.endsWith("/") ? configuredUrl : `${configuredUrl}/`);
}
