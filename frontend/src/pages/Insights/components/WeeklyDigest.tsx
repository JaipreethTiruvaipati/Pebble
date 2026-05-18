import { weeklyDigest } from "@/lib/firebase";
import { formatCurrency } from "@/lib/formatCurrency";
import { TrendingDown } from "lucide-react";

export function WeeklyDigest() {
  const diff = weeklyDigest.thisWeek - weeklyDigest.lastWeek;
  const pct = ((diff / weeklyDigest.lastWeek) * 100).toFixed(1);
  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <div className="flex items-center justify-between">
        <h3 className="text-lg">Weekly digest</h3>
        <span className="text-[11px] text-muted-foreground">Mon · Sun</span>
      </div>
      <div className="mt-5 grid grid-cols-2 gap-4">
        <div className="rounded-2xl bg-card-elevated p-4">
          <div className="text-[11px] uppercase tracking-widest text-muted-foreground">This week</div>
          <div className="mt-2 font-mono text-2xl">{formatCurrency(weeklyDigest.thisWeek)}</div>
        </div>
        <div className="rounded-2xl bg-card-elevated p-4">
          <div className="text-[11px] uppercase tracking-widest text-muted-foreground">Last week</div>
          <div className="mt-2 font-mono text-2xl text-muted-foreground">{formatCurrency(weeklyDigest.lastWeek)}</div>
        </div>
      </div>
      <div className="mt-4 flex items-center gap-2 rounded-xl bg-teal/10 px-4 py-3 text-sm text-teal">
        <TrendingDown size={14} /> {pct}% vs last week — saved {formatCurrency(weeklyDigest.saved)}
      </div>
      <div className="mt-2 text-xs text-muted-foreground">Top category: <span className="text-foreground">{weeklyDigest.topCategory}</span></div>
    </div>
  );
}
export default WeeklyDigest;

