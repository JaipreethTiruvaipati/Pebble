import { motion } from "framer-motion";
import { categoryBreakdown } from "@/lib/firebase";
import { formatCurrency } from "@/lib/formatCurrency";
import { scoreColor } from "@/lib/scoreColor";

export function CategoryBarChart() {
  const max = Math.max(...categoryBreakdown.map((c) => c.spend));
  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <h3 className="text-lg">Spend by category</h3>
      <p className="mt-1 text-xs text-muted-foreground">Bar color = average impulse score for that category.</p>
      <ul className="mt-6 space-y-4">
        {categoryBreakdown.map((c, i) => {
          const pct = (c.spend / max) * 100;
          return (
            <li key={c.category}>
              <div className="flex items-center justify-between text-sm">
                <span>{c.category}</span>
                <span className="font-mono">{formatCurrency(c.spend)}</span>
              </div>
              <div className="mt-1.5 h-3 overflow-hidden rounded-full bg-card-elevated">
                <motion.div
                  initial={{ width: 0 }} animate={{ width: `${pct}%` }}
                  transition={{ delay: i * 0.08, duration: 0.9, ease: [0.16, 1, 0.3, 1] }}
                  className="h-full rounded-full"
                  style={{ background: scoreColor(c.score) }}
                />
              </div>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
export default CategoryBarChart;

