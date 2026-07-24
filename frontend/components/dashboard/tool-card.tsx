import Link from "next/link";
import { ArrowUpRight } from "lucide-react";

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
      className="group relative flex min-h-[17rem] flex-col overflow-hidden rounded-lg border border-line bg-surface p-6 transition-[border-color,box-shadow] duration-200 hover:border-line-strong hover:shadow-soft"
    >
      <div className="flex items-start justify-between">
        <span className="font-mono text-xs font-semibold tracking-[0.12em] text-ink-muted">
          {String(index + 1).padStart(2, "0")}
        </span>
        <span className="grid size-10 place-items-center rounded-md bg-canvas text-ink-muted transition-colors group-hover:bg-accent group-hover:text-ink">
          <Icon aria-hidden="true" className="size-5" />
        </span>
      </div>

      <div className="mt-auto pt-12">
        <p className="font-mono text-[0.65rem] font-semibold uppercase tracking-[0.14em] text-ink-muted">
          {tool.accepts} → {tool.output}
        </p>
        <h3 className="mt-3 text-xl leading-tight font-semibold tracking-[-0.025em]">
          {tool.title}
        </h3>
        <p className="mt-3 max-w-sm text-sm leading-6 text-ink-muted">
          {tool.description}
        </p>
      </div>

      <div className="mt-6 flex items-center justify-between border-t border-line pt-4 text-xs font-semibold text-ink-muted">
        Open tool
        <ArrowUpRight
          aria-hidden="true"
          className="size-4 transition-transform group-hover:translate-x-1 group-hover:-translate-y-1"
        />
      </div>
    </Link>
  );
}
