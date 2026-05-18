import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { CategoryBarChart } from "./components/CategoryBarChart";
import { WeeklyDigest } from "./components/WeeklyDigest";
import { PeerBenchmark } from "./components/PeerBenchmark";
import { StreakCounter } from "./components/StreakCounter";
import { topMerchants } from "@/lib/firebase";
import { formatCurrency } from "@/lib/formatCurrency";
import { scoreBgClass } from "@/lib/scoreColor";

function TopMerchants() {
  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <h3 className="text-lg">Top merchants</h3>
      <ol className="mt-4 space-y-3">
        {topMerchants.map((m, i) => (
          <li key={m.merchant} className="flex items-center gap-4">
            <span className="w-6 font-mono text-sm text-muted-foreground">{String(i + 1).padStart(2, "0")}</span>
            <div className="flex-1">
              <div className="text-sm">{m.merchant}</div>
            </div>
            <span className="font-mono text-sm">{formatCurrency(m.spend)}</span>
            <span className={`rounded-full px-2 py-0.5 font-mono text-[11px] ${scoreBgClass(m.score)}`}>{m.score}</span>
          </li>
        ))}
      </ol>
    </div>
  );
}

export default function Insights() {
  return (
    <AppShell>
      <div className="mx-auto max-w-7xl p-6 md:p-10">
        <PageHeader title="Insights" subtitle="Patterns in how you spend, weekly." />
        <div className="grid gap-6 lg:grid-cols-[1.5fr_1fr]">
          <CategoryBarChart />
          <StreakCounter />
        </div>
        <div className="mt-6 grid gap-6 lg:grid-cols-3">
          <WeeklyDigest />
          <TopMerchants />
          <PeerBenchmark />
        </div>
      </div>
    </AppShell>
  );
}

