import { motion } from "framer-motion";
import { formatCurrency } from "@/lib/formatCurrency";
import { scoreColor } from "@/lib/scoreColor";
import { useWeeklyInsights } from "@/hooks/usePortfolio";
import { Skeleton } from "@/components/ui/Skeleton";

export function CategoryBarChart() {
  const { data, isLoading } = useWeeklyInsights();
  if (isLoading || !data) return <Skeleton className="h-80 rounded-3xl" />;

  const categories = data.top_categories.length
    ? data.top_categories
    : [{ category: "No data", amount: 0, pct: 0 }];
  const max = Math.max(...categories.map((c) => c.amount), 1);

  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <h3 className="text-lg">Spend by category</h3>
      <p className="mt-1 text-xs text-muted-foreground">From your last 7 days of transactions.</p>
      <ul className="mt-6 space-y-4">
        {categories.map((c, i) => {
          const pct = (c.amount / max) * 100;
          const score = data.avg_impulse_score;
          return (
            <li key={c.category}>
              <div className="flex items-center justify-between text-sm">
                <span>{c.category}</span>
                <span className="font-mono">{formatCurrency(c.amount)}</span>
              </div>
              <div className="mt-1.5 h-3 overflow-hidden rounded-full bg-card-elevated">
                <motion.div
                  initial={{ width: 0 }}
                  animate={{ width: `${pct}%` }}
                  transition={{ delay: i * 0.08, duration: 0.9, ease: [0.16, 1, 0.3, 1] }}
                  className="h-full rounded-full"
                  style={{ background: scoreColor(score) }}
                />
              </div>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
