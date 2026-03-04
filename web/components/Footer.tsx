import { Github, Linkedin, Twitter } from "lucide-react";

const links = [
  { name: "GitHub", href: "https://github.com/saral-gupta7", icon: Github },
  { name: "Twitter", href: "https://x.com/srlgpt_a", icon: Twitter },
  {
    name: "LinkedIn",
    href: "https://linkedin.com/saralgupta7",
    icon: Linkedin,
  },
];

export default function Footer() {
  return (
    <footer className="mt-10 border-t border-zinc-800 pt-6">
      <div className="flex flex-col gap-4 text-sm text-zinc-400 sm:flex-row items-center sm:justify-between justify-center">
        <p>Built with Next.js and raw ffmpeg.</p>
        <div className="flex items-center gap-3">
          {links.map((link) => {
            const Icon = link.icon;
            return (
              <a
                key={link.name}
                href={link.href}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center gap-2 rounded-lg border border-zinc-700 px-3 py-2 text-zinc-300 transition hover:border-zinc-500 hover:text-white"
              >
                <Icon size={15} />
                {link.name}
              </a>
            );
          })}
        </div>
      </div>
    </footer>
  );
}
