import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { useTransactions } from "@/hooks/useTransaction";
import { formatCurrency } from "@/lib/formatCurrency";
import { formatRelative } from "@/lib/formatDate";
import { scoreBgClass } from "@/lib/scoreColor";
import { Skeleton } from "@/components/ui/Skeleton";

export default function TransactionHistory() {
  const { data, isLoading } = useTransactions(50);

  return (
    <AppShell>
      <div className="mx-auto max-w-4xl p-6 md:p-10">
        <PageHeader title="Transaction history" subtitle="All logged spends and impulse scores." />
        {isLoading ? (
          <Skeleton className="h-64" />
        ) : !data?.length ? (
          <p className="text-muted-foreground">No transactions yet.</p>
        ) : (
          <ul className="divide-y divide-border rounded-3xl border border-border bg-card">
            {data.map((t) => (
              <li key={t.id} className="flex items-center gap-4 px-5 py-4">
                <div className="grid h-10 w-10 place-items-center rounded-xl bg-card-elevated text-xs font-medium">
                  {t.merchant.slice(0, 2).toUpperCase()}
                </div>
                <div className="min-w-0 flex-1">
                  <div className="text-sm">{t.merchant}</div>
                  <div className="text-[11px] text-muted-foreground">
                    {formatRelative(t.logged_at)} · {t.status}
                  </div>
                </div>
                <div className="text-right">
                  <div className="font-mono text-sm">{formatCurrency(t.total_amount)}</div>
                  {t.total_penalty > 0 && (
                    <div className="font-mono text-[11px] text-coral">−{formatCurrency(t.total_penalty)}</div>
                  )}
                </div>
                <span className={`rounded-full px-2 py-0.5 font-mono text-xs ${scoreBgClass(Math.round(t.avg_score))}`}>
                  {Math.round(t.avg_score)}
                </span>
              </li>
            ))}
          </ul>
        )}
      </div>
    </AppShell>
  );
}
