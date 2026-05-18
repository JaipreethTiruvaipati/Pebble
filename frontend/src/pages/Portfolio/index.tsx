import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { MetricCard } from "@/components/ui/MetricCard";
import { usePortfolio, useInvestments } from "@/hooks/usePortfolio";
import { formatCurrency } from "@/lib/formatCurrency";
import { Skeleton } from "@/components/ui/Skeleton";
import { PieChart, Pie, Cell, ResponsiveContainer } from "recharts";

const COLORS = ["#E94560", "#F5A623", "#00D4AA", "#8B8FA8"];

export default function Portfolio() {
  const { data: portfolio, isLoading } = usePortfolio();
  const { data: invList } = useInvestments(20);

  const total =
    (portfolio?.equity_value ?? 0) + (portfolio?.gold_value ?? 0) + (portfolio?.bond_value ?? 0);
  const allocation = portfolio?.allocation_pct
    ? Object.entries(portfolio.allocation_pct).map(([name, value]) => ({ name, value }))
    : [];

  if (isLoading) {
    return (
      <AppShell>
        <Skeleton className="m-10 h-96" />
      </AppShell>
    );
  }

  return (
    <AppShell>
      <div className="mx-auto max-w-7xl p-6 md:p-10">
        <PageHeader title="Portfolio" subtitle="Where your penalties became investments." />

        <div className="grid gap-4 md:grid-cols-3">
          <MetricCard label="Total invested" value={portfolio?.total_invested ?? 0} currency />
          <MetricCard label="Current value" value={total} currency accent="teal" />
          <MetricCard label="Returns" value={portfolio?.gain_pct ?? 0} accent="teal" footer="gain %" />
        </div>

        <div className="mt-6 grid gap-6 lg:grid-cols-2">
          <div className="rounded-3xl border border-border bg-card p-6">
            <h3 className="text-lg">Allocation</h3>
            <div className="mt-4 h-48">
              <ResponsiveContainer>
                <PieChart>
                  <Pie data={allocation} dataKey="value" innerRadius={50} outerRadius={80} strokeWidth={0}>
                    {allocation.map((_, i) => (
                      <Cell key={i} fill={COLORS[i % COLORS.length]} />
                    ))}
                  </Pie>
                </PieChart>
              </ResponsiveContainer>
            </div>
          </div>
          <div className="rounded-3xl border border-border bg-card p-6">
            <h3 className="text-lg">Recent investments</h3>
            <ul className="mt-4 divide-y divide-border">
              {(invList?.investments ?? []).map((inv) => (
                <li key={inv.id} className="flex justify-between py-3 text-sm">
                  <span>
                    {inv.asset_class} · {inv.trigger_type || "batch"}
                  </span>
                  <span className="font-mono">{formatCurrency(inv.amount)}</span>
                </li>
              ))}
              {!invList?.investments?.length && (
                <li className="py-4 text-muted-foreground">No investments yet — confirm penalties to fund the pool.</li>
              )}
            </ul>
          </div>
        </div>
      </div>
    </AppShell>
  );
}
