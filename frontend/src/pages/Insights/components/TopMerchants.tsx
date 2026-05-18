import { useMemo } from "react";
import { useTransactions } from "@/hooks/useTransaction";
import { formatCurrency } from "@/lib/formatCurrency";
import { scoreBgClass } from "@/lib/scoreColor";
import { Skeleton } from "@/components/ui/Skeleton";

export function TopMerchants() {
  const { data: txs, isLoading } = useTransactions(50);

  const merchants = useMemo(() => {
    const map = new Map<string, { merchant: string; spend: number; scoreSum: number; count: number }>();
    for (const tx of txs ?? []) {
      const cur = map.get(tx.merchant) ?? { merchant: tx.merchant, spend: 0, scoreSum: 0, count: 0 };
      cur.spend += tx.total_amount;
      cur.scoreSum += tx.avg_score ?? 0;
      cur.count += 1;
      map.set(tx.merchant, cur);
    }
    return [...map.values()]
      .map((m) => ({
        merchant: m.merchant,
        spend: m.spend,
        score: m.count ? Math.round(m.scoreSum / m.count) : 0,
      }))
      .sort((a, b) => b.spend - a.spend)
      .slice(0, 5);
  }, [txs]);

  if (isLoading) return <Skeleton className="h-56 rounded-3xl" />;

  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <h3 className="text-lg">Top merchants</h3>
      {merchants.length === 0 ? (
        <p className="mt-4 text-sm text-muted-foreground">Log transactions to see merchant rankings.</p>
      ) : (
        <ol className="mt-4 space-y-3">
          {merchants.map((m, i) => (
            <li key={m.merchant} className="flex items-center gap-4">
              <span className="w-6 font-mono text-sm text-muted-foreground">{String(i + 1).padStart(2, "0")}</span>
              <div className="flex-1">
                <div className="text-sm">{m.merchant}</div>
              </div>
              <span className="font-mono text-sm">{formatCurrency(m.spend)}</span>
              <span className={`rounded-full px-2 py-0.5 font-mono text-[11px] ${scoreBgClass(m.score)}`}>
                {m.score}
              </span>
            </li>
          ))}
        </ol>
      )}
    </div>
  );
}
