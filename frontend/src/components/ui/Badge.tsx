import { cn } from "@/lib/cn";
import type { ReactNode } from "react";

type Variant = "default" | "coral" | "teal" | "amber" | "outline";
const variants: Record<Variant, string> = {
  default: "bg-card-elevated text-foreground",
  coral: "bg-coral/15 text-coral",
  teal: "bg-teal/15 text-teal",
  amber: "bg-amber/15 text-amber",
  outline: "border border-border text-muted-foreground",
};

export function Badge({ children, variant = "default", className }: { children: ReactNode; variant?: Variant; className?: string }) {
  return (
    <span className={cn("inline-flex items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-medium", variants[variant], className)}>
      {children}
    </span>
  );
}
export default Badge;

