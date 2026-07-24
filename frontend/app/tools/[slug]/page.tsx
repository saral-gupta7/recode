import type { Metadata } from "next";
import { notFound } from "next/navigation";

import { MediaJobWorkspace } from "@/components/tools/media-job-workspace";
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
    return {};
  }

  return {
    title: tool.title,
    description: tool.description,
  };
}

export default async function ToolPage({ params }: ToolPageProps) {
  const { slug } = await params;
  const tool = getToolBySlug(slug);

  if (!tool) {
    notFound();
  }

  return <MediaJobWorkspace toolSlug={tool.slug} />;
}
