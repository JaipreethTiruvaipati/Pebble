import { formatCurrency } from "@/lib/formatCurrency";
import { TrendingDown, TrendingUp } from "lucide-react";
import { useWeeklyInsights } from "@/hooks/usePortfolio";
import { Skeleton } from "@/components/ui/Skeleton";

export function WeeklyDigest() {
  const { data, isLoading } = useWeeklyInsights();
  if (isLoading || !data) return <Skeleton className="h-56 rounded-3xl" />;

  const diff = data.trend_vs_last_week_pct;
  const topCat = data.top_categories?.[0]?.category ?? "—";

  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <div className="flex items-center justify-between">
        <h3 className="text-lg">Weekly digest</h3>
        <span className="text-[11px] text-muted-foreground">Last 7 days</span>
      </div>
      <div className="mt-5 grid grid-cols-2 gap-4">
        <div className="rounded-2xl bg-card-elevated p-4">
          <div className="text-[11px] uppercase tracking-widest text-muted-foreground">This week</div>
          <div className="mt-2 font-mono text-2xl">{formatCurrency(data.total_spend)}</div>
        </div>
        <div className="rounded-2xl bg-card-elevated p-4">
          <div className="text-[11px] uppercase tracking-widest text-muted-foreground">Impulse %</div>
          <div className="mt-2 font-mono text-2xl">{data.impulse_pct.toFixed(0)}%</div>
        </div>
      </div>
      <div
        className={`mt-4 flex items-center gap-2 rounded-xl px-4 py-3 text-sm ${
          diff <= 0 ? "bg-teal/10 text-teal" : "bg-coral/10 text-coral"
        }`}
      >
        {diff <= 0 ? <TrendingDown size={14} /> : <TrendingUp size={14} />}
        {Math.abs(diff).toFixed(1)}% vs prior week
      </div>
      <div className="mt-2 text-xs text-muted-foreground">
        Top category: <span className="text-foreground">{topCat}</span>
      </div>
    </div>
  );
}
