import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { CategoryBarChart } from "./components/CategoryBarChart";
import { WeeklyDigest } from "./components/WeeklyDigest";
import { PeerBenchmark } from "./components/PeerBenchmark";
import { StreakCounter } from "./components/StreakCounter";
import { TopMerchants } from "./components/TopMerchants";

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

