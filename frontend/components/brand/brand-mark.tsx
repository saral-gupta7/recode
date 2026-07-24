import Link from "next/link";

interface BrandMarkProps {
  compact?: boolean;
}

export function BrandMark({ compact = false }: BrandMarkProps) {
  return (
    <Link
      href="/"
      aria-label="Recode home"
      className="group inline-flex items-center gap-3"
    >
      <span
        aria-hidden="true"
        className="grid size-10 place-items-center border-2 border-ink bg-accent font-mono text-xs font-semibold text-[#171a17] shadow-hard transition-transform duration-200 group-hover:-rotate-3"
      >
        R/
      </span>
      {!compact && (
        <span className="text-xl font-bold uppercase tracking-[-0.04em]">
          Recode
        </span>
      )}
    </Link>
  );
}
