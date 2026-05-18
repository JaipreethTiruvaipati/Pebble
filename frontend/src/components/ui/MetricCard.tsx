import { motion, useMotionValue, useTransform, animate } from "framer-motion";
import { useEffect, type ReactNode } from "react";
import { cn } from "@/lib/cn";
import { formatCurrency, formatNumber } from "@/lib/formatCurrency";

export function CountUp({ value, currency = false, compact = false, className, duration = 1.2 }: { value: number; currency?: boolean; compact?: boolean; className?: string; duration?: number }) {
  const mv = useMotionValue(0);
  const rounded = useTransform(mv, (v) => (currency ? formatCurrency(Math.round(v), { compact }) : formatNumber(Math.round(v))));
  useEffect(() => {
    const ctrl = animate(mv, value, { duration, ease: [0.16, 1, 0.3, 1] });
    return ctrl.stop;
  }, [value, duration, mv]);
  return <motion.span className={cn("font-mono", className)}>{rounded}</motion.span>;
}

export function MetricCard({ label, value, delta, icon, accent = "default", currency = false, compact = false, footer }: {
  label: string;
  value: number;
  delta?: { value: number; positive?: boolean; suffix?: string };
  icon?: ReactNode;
  accent?: "default" | "coral" | "teal" | "amber";
  currency?: boolean;
  compact?: boolean;
  footer?: ReactNode;
}) {
  const accentClass = {
    default: "text-foreground",
    coral: "text-coral",
    teal: "text-teal",
    amber: "text-amber",
  }[accent];
  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5, ease: [0.16, 1, 0.3, 1] }}
      className="relative overflow-hidden rounded-2xl bg-card border border-border p-5"
    >
      <div className="flex items-start justify-between gap-2">
        <span className="text-xs uppercase tracking-wider text-muted-foreground">{label}</span>
        {icon && <span className="text-muted-foreground">{icon}</span>}
      </div>
      <div className={cn("mt-3 text-3xl font-mono font-semibold tracking-tight", accentClass)}>
        <CountUp value={value} currency={currency} compact={compact} />
      </div>
      {delta && (
        <div className={cn("mt-1 text-xs font-mono", delta.positive ? "text-teal" : "text-coral")}>
          {delta.positive ? "▲" : "▼"} {Math.abs(delta.value)}{delta.suffix ?? "%"}
        </div>
      )}
      {footer && <div className="mt-3 text-xs text-muted-foreground">{footer}</div>}
    </motion.div>
  );
}
export default MetricCard;

