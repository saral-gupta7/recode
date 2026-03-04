"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Layers, Music, Scissors, VolumeX, Wand2, Video } from "lucide-react";

const tools = [
  { name: "Grayscale", href: "/grayscale", icon: Wand2, color: "text-zinc-300" },
  { name: "Extract Audio", href: "/extract-audio", icon: Music, color: "text-emerald-400" },
  { name: "Remove Audio", href: "/remove-audio", icon: VolumeX, color: "text-cyan-400" },
  { name: "Trim 30s", href: "/trim", icon: Scissors, color: "text-violet-400" },
  { name: "Compress", href: "/compress", icon: Video, color: "text-orange-400" },
  { name: "Frame Extract", href: "/frames", icon: Scissors, color: "text-amber-400" },
];

function NavLinks({ compact = false }: { compact?: boolean }) {
  const pathname = usePathname();

  return (
    <nav className={compact ? "flex gap-2 overflow-x-auto pb-1" : "flex flex-col gap-1.5"}>
      {!compact && (
        <p className="mb-2 px-4 text-[10px] font-bold uppercase tracking-widest text-zinc-500">
          Media Tools
        </p>
      )}
      {tools.map((tool) => {
        const Icon = tool.icon;
        const isActive = pathname === tool.href;

        return (
          <Link
            key={tool.href}
            href={tool.href}
            className={`group flex items-center gap-3 rounded-xl transition-all ${
              compact ? "px-3 py-2 whitespace-nowrap" : "px-4 py-3"
            } ${
              isActive
                ? "bg-zinc-800 text-white"
                : "text-zinc-400 hover:bg-zinc-800/50 hover:text-zinc-200"
            }`}
          >
            <Icon
              size={17}
              className={isActive ? tool.color : "text-zinc-500 group-hover:text-zinc-300"}
            />
            <span className="text-sm font-medium">{tool.name}</span>
          </Link>
        );
      })}
    </nav>
  );
}

export default function Sidebar() {
  return (
    <>
      <div className="border-b border-zinc-800 bg-zinc-900/90 p-4 backdrop-blur md:hidden">
        <div className="mb-3 flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-cyan-500 to-blue-600">
            <Layers size={17} className="text-white" />
          </div>
          <span className="text-base font-semibold text-white">Recode</span>
        </div>
        <NavLinks compact />
      </div>

      <aside className="sticky top-0 hidden h-screen w-72 border-r border-zinc-800 bg-zinc-900/50 p-6 md:flex md:flex-col">
        <div className="mb-8 flex items-center gap-3 px-2">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-gradient-to-br from-cyan-500 to-blue-600">
            <Layers size={20} className="text-white" />
          </div>
          <span className="text-xl font-bold tracking-tight text-white">Recode</span>
        </div>
        <NavLinks />
      </aside>
    </>
  );
}
