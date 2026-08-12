import Link from "next/link";
import { ArrowUpRight } from "@phosphor-icons/react/ssr";

import type { MediaTool } from "@/config/tools";

interface ToolCardProps {
  index: number;
  tool: MediaTool;
}

export function ToolCard({ index, tool }: ToolCardProps) {
  const Icon = tool.icon;

  return (
    <Link
      href={`/tools/${tool.slug}`}
      className="group grid min-h-32 grid-cols-[2.5rem_minmax(0,1fr)_1.5rem] items-center gap-3 rounded-md border border-line bg-surface p-4 shadow-hard hover:-translate-y-0.5 hover:border-line-strong hover:shadow-soft"
    >
      <span className="grid size-10 place-items-center rounded-md bg-canvas text-ink-muted group-hover:bg-accent group-hover:text-accent-strong">
        <Icon aria-hidden="true" weight="bold" className="size-4.5" />
      </span>

      <div className="min-w-0">
        <div className="flex items-center gap-2">
          <span className="font-mono text-[0.6rem] font-semibold text-ink-muted">{String(index + 1).padStart(2, "0")}</span>
          <span className="truncate text-[0.62rem] font-medium text-ink-muted">{tool.output}</span>
        </div>
        <h3 className="mt-1 text-[0.95rem] leading-tight font-semibold tracking-[-0.02em]">
          {tool.title}
        </h3>
        <p className="mt-1.5 line-clamp-2 text-xs leading-5 text-ink-muted">
          {tool.description}
        </p>
      </div>

      <ArrowUpRight aria-hidden="true" className="size-4 text-ink-muted group-hover:translate-x-0.5 group-hover:-translate-y-0.5 group-hover:text-accent-strong" />
    </Link>
  );
}
