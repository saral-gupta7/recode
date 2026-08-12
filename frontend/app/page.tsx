import type { Metadata } from "next";

import { ActiveJobBanner } from "@/components/dashboard/active-job-banner";
import { DashboardHero } from "@/components/dashboard/dashboard-hero";
import { HowItWorks } from "@/components/dashboard/how-it-works";
import { ToolCatalog } from "@/components/dashboard/tool-catalog";
import { JsonLd } from "@/components/seo/json-ld";
import { getSiteUrl, siteConfig } from "@/config/site";
import { mediaTools } from "@/config/tools";

export const metadata: Metadata = {
  title: { absolute: siteConfig.title },
  description: siteConfig.description,
  alternates: { canonical: "/" },
};

export default function HomePage() {
  const siteUrl = getSiteUrl();
  const toolList = mediaTools.map((tool, index) => ({
    "@type": "ListItem",
    position: index + 1,
    name: tool.title,
    url: new URL(`/tools/${tool.slug}`, siteUrl).toString(),
  }));

  return (
    <>
      <JsonLd
        data={[
          {
            "@context": "https://schema.org",
            "@type": "WebApplication",
            name: siteConfig.name,
            url: siteUrl.toString(),
            description: siteConfig.description,
            applicationCategory: "MultimediaApplication",
            operatingSystem: "Any",
            browserRequirements: "Requires JavaScript and a modern web browser",
            offers: {
              "@type": "Offer",
              price: "0",
              priceCurrency: "USD",
            },
            featureList: mediaTools.map((tool) => tool.title),
          },
          {
            "@context": "https://schema.org",
            "@type": "ItemList",
            name: "Recode image tools",
            numberOfItems: mediaTools.length,
            itemListElement: toolList,
          },
        ]}
      />
      <main id="main-content" className="mx-auto w-full max-w-6xl px-4 pb-12 sm:px-6 lg:px-8">
        <ActiveJobBanner />
        <DashboardHero />
        <ToolCatalog />
        <HowItWorks />
      </main>
    </>
  );
}
