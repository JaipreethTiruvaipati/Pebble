import { Link } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { formatCurrency } from "@/lib/formatCurrency";
import { formatRelative } from "@/lib/formatDate";
import { scoreBgClass } from "@/lib/scoreColor";
import { useTransactions } from "@/hooks/useTransaction";
import { ROUTES } from "@/routes";
import { Skeleton } from "@/components/ui/Skeleton";

export function RecentTransactions() {
  const { data, isLoading } = useTransactions(6);

  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <div className="flex items-center justify-between">
        <h3 className="text-lg">Recent transactions</h3>
        <Link to={ROUTES.history} className="text-xs text-muted-foreground hover:text-foreground">
          View all →
        </Link>
      </div>
      {isLoading ? (
        <Skeleton className="mt-4 h-48" />
      ) : !data?.length ? (
        <p className="mt-4 text-sm text-muted-foreground">No transactions yet. Log your first bill.</p>
      ) : (
        <ul className="mt-4 divide-y divide-border">
          {data.map((t, i) => (
            <motion.li
              key={t.id}
              initial={{ opacity: 0, x: -8 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{ delay: i * 0.06, duration: 0.4 }}
              className="flex items-center gap-4 py-3"
            >
              <div className="grid h-10 w-10 place-items-center rounded-xl bg-card-elevated text-xs font-medium">
                {t.merchant.slice(0, 2).toUpperCase()}
              </div>
              <div className="min-w-0 flex-1">
                <div className="truncate text-sm">{t.merchant}</div>
                <div className="mt-0.5 text-[11px] text-muted-foreground">{formatRelative(t.logged_at)}</div>
              </div>
              <div className="text-right">
                <div className="font-mono text-sm">{formatCurrency(t.total_amount)}</div>
                {t.total_penalty > 0 && (
                  <div className="font-mono text-[11px] text-coral">−{formatCurrency(t.total_penalty)}</div>
                )}
              </div>
              <span
                className={`ml-2 rounded-full px-2 py-0.5 font-mono text-xs ${scoreBgClass(Math.round(t.avg_score))}`}
              >
                {Math.round(t.avg_score)}
              </span>
            </motion.li>
          ))}
        </ul>
      )}
    </div>
  );
}
