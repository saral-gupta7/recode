import type { Metadata } from "next";
import { notFound } from "next/navigation";

import { MediaJobWorkspace } from "@/components/tools/media-job-workspace";
import { JsonLd } from "@/components/seo/json-ld";
import { getSiteUrl, siteConfig } from "@/config/site";
import { getToolBySlug, mediaTools } from "@/config/tools";

interface ToolPageProps {
  params: Promise<{ slug: string }>;
}

export function generateStaticParams() {
  return mediaTools.map((tool) => ({ slug: tool.slug }));
}

export async function generateMetadata({
  params,
}: ToolPageProps): Promise<Metadata> {
  const { slug } = await params;
  const tool = getToolBySlug(slug);

  if (!tool) {
    return { robots: { index: false, follow: false } };
  }

  const path = `/tools/${tool.slug}`;
  const description = `${tool.description} Free to use with no account or advertising.`;

  return {
    title: tool.title,
    description,
    alternates: { canonical: path },
    openGraph: {
      type: "website",
      url: path,
      siteName: siteConfig.name,
      title: `${tool.title} · Recode`,
      description,
      images: ["/opengraph-image"],
    },
    twitter: {
      card: "summary_large_image",
      title: `${tool.title} · Recode`,
      description,
      images: ["/opengraph-image"],
    },
  };
}

export default async function ToolPage({ params }: ToolPageProps) {
  const { slug } = await params;
  const tool = getToolBySlug(slug);

  if (!tool) {
    notFound();
  }

  const url = new URL(`/tools/${tool.slug}`, getSiteUrl()).toString();

  return (
    <>
      <JsonLd
        data={[
          {
            "@context": "https://schema.org",
            "@type": "WebApplication",
            name: `${tool.title} — Recode`,
            url,
            description: tool.description,
            applicationCategory: "MultimediaApplication",
            operatingSystem: "Any",
            isPartOf: {
              "@type": "WebSite",
              name: siteConfig.name,
              url: getSiteUrl().toString(),
            },
            offers: {
              "@type": "Offer",
              price: "0",
              priceCurrency: "USD",
            },
          },
          {
            "@context": "https://schema.org",
            "@type": "BreadcrumbList",
            itemListElement: [
              {
                "@type": "ListItem",
                position: 1,
                name: "Image tools",
                item: getSiteUrl().toString(),
              },
              {
                "@type": "ListItem",
                position: 2,
                name: tool.title,
                item: url,
              },
            ],
          },
        ]}
      />
      <MediaJobWorkspace toolSlug={tool.slug} />
    </>
  );
}
