import { motion } from "framer-motion";
import { CategoryChip } from "@/components/ui/CategoryChip";
import { recentTransactions } from "@/lib/firebase";
import { formatCurrency } from "@/lib/formatCurrency";
import { formatRelative } from "@/lib/formatDate";
import { scoreBgClass } from "@/lib/scoreColor";

export function RecentTransactions() {
  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <div className="flex items-center justify-between">
        <h3 className="text-lg">Recent transactions</h3>
        <button className="text-xs text-muted-foreground hover:text-foreground">View all →</button>
      </div>
      <ul className="mt-4 divide-y divide-border">
        {recentTransactions.map((t, i) => (
          <motion.li
            key={t.id}
            initial={{ opacity: 0, x: -8 }} animate={{ opacity: 1, x: 0 }}
            transition={{ delay: i * 0.06, duration: 0.4 }}
            className="flex items-center gap-4 py-3"
          >
            <div className="grid h-10 w-10 place-items-center rounded-xl bg-card-elevated text-xs font-medium">
              {t.merchant.slice(0, 2).toUpperCase()}
            </div>
            <div className="min-w-0 flex-1">
              <div className="truncate text-sm">{t.merchant}</div>
              <div className="mt-0.5 flex items-center gap-2 text-[11px] text-muted-foreground">
                <CategoryChip category={t.category} />
                <span>{formatRelative(t.date)}</span>
              </div>
            </div>
            <div className="text-right">
              <div className="font-mono text-sm">{formatCurrency(t.amount)}</div>
              {t.penalty > 0 && <div className="font-mono text-[11px] text-coral">−{formatCurrency(t.penalty)}</div>}
            </div>
            <span className={`ml-2 rounded-full px-2 py-0.5 font-mono text-xs ${scoreBgClass(t.score)}`}>{t.score}</span>
          </motion.li>
        ))}
      </ul>
    </div>
  );
}
export default RecentTransactions;

