"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { Video, Music, Wand2, Scissors, Layers } from "lucide-react";

const tools = [
  {
    name: "Video to GIF",
    href: "/",
    icon: Video,
    color: "text-purple-400",
  },
  {
    name: "Extract Audio",
    href: "/extract-audio",
    icon: Music,
    color: "text-blue-400",
  },
  {
    name: "Black & White",
    href: "/grayscale",
    icon: Wand2,
    color: "text-zinc-400",
  },
  {
    name: "Frame Extractor",
    href: "/frames",
    icon: Scissors,
    color: "text-emerald-400",
  },
];

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="w-64 border-r border-zinc-800 bg-zinc-900/50 p-6 flex flex-col gap-8 h-screen sticky top-0">
      <div className="flex items-center gap-3 px-2">
        <div className="h-9 w-9 rounded-xl bg-gradient-to-br from-blue-600 to-indigo-700 flex items-center justify-center shadow-lg shadow-blue-900/20">
          <Layers size={20} className="text-white" />
        </div>
        <span className="text-xl font-bold tracking-tight text-white">
          Recode
        </span>
      </div>

      <nav className="flex flex-col gap-1.5">
        <p className="text-[10px] font-bold uppercase tracking-widest text-zinc-500 px-4 mb-2">
          Media Tools
        </p>
        {tools.map((tool) => {
          const Icon = tool.icon;
          const isActive = pathname === tool.href;
          return (
            <Link
              key={tool.href}
              href={tool.href}
              className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 group ${
                isActive
                  ? "bg-zinc-800 text-white shadow-sm"
                  : "text-zinc-400 hover:bg-zinc-800/50 hover:text-zinc-200"
              }`}
            >
              <Icon
                size={18}
                className={`${isActive ? tool.color : "text-zinc-500 group-hover:text-zinc-300"}`}
              />
              <span className="text-sm font-medium">{tool.name}</span>
              {isActive && (
                <div className="ml-auto h-1.5 w-1.5 rounded-full bg-blue-500" />
              )}
            </Link>
          );
        })}
      </nav>

      <div className="mt-auto border-t border-zinc-800 pt-6 px-2"></div>
    </aside>
  );
}
