import { cn } from "@/lib/cn";
import type { ReactNode } from "react";

const categoryColors: Record<string, string> = {
  Food: "bg-amber/15 text-amber",
  Shopping: "bg-coral/15 text-coral",
  Entertainment: "bg-coral/15 text-coral",
  Travel: "bg-teal/15 text-teal",
  Bills: "bg-muted text-muted-foreground",
  Groceries: "bg-teal/15 text-teal",
  Tech: "bg-coral/15 text-coral",
  Health: "bg-teal/15 text-teal",
  Other: "bg-muted text-muted-foreground",
};

export function CategoryChip({ category, className }: { category: string; className?: string }) {
  const color = categoryColors[category] ?? "bg-muted text-muted-foreground";
  return (
    <span className={cn("inline-flex items-center rounded-md px-2 py-0.5 text-[11px] font-medium", color, className)}>
      {category}
    </span>
  );
}

export function ChipRow({ children }: { children: ReactNode }) {
  return <div className="flex flex-wrap gap-2">{children}</div>;
}
export default CategoryChip;

