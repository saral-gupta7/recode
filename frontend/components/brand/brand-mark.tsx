import Link from "next/link";

interface BrandMarkProps {
  compact?: boolean;
}

export function BrandMark({ compact = false }: BrandMarkProps) {
  return (
    <Link
      href="/"
      aria-label="Recode home"
      className="group inline-flex items-center gap-2.5"
    >
      <span
        aria-hidden="true"
        className="grid size-9 place-items-center rounded-md bg-ink font-mono text-[0.68rem] font-semibold text-canvas shadow-hard group-hover:-rotate-3"
      >
        R/
      </span>
      {!compact && (
        <span className="text-lg font-bold tracking-[-0.04em]">
          Recode
        </span>
      )}
    </Link>
  );
}
