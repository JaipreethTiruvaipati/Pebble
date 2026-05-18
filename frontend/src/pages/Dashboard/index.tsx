import { Link } from "@tanstack/react-router";
import { Upload } from "lucide-react";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { WalletCard } from "./components/WalletCard";
import { RecentTransactions } from "./components/RecentTransactions";
import { PendingPenaltyBanner } from "./components/PendingPenaltyBanner";
import { PortfolioMiniChart } from "./components/PortfolioMiniChart";
import { MarketSignalIndicator } from "./components/MarketSignalIndicator";
import { ImpulseStreakBadge } from "./components/ImpulseStreakBadge";
import { ROUTES } from "@/routes";
import { useAuthStore } from "@/stores/authStore";

export default function Dashboard() {
  const { user } = useAuthStore();
  return (
    <AppShell>
      <div className="mx-auto max-w-7xl p-6 md:p-10">
        <PageHeader
          title={`Hi, ${user?.name.split(" ")[0]}.`}
          subtitle="Here's where your money sits today."
          actions={<ImpulseStreakBadge />}
        />

        <PendingPenaltyBanner />

        <div className="mt-6 grid gap-6 lg:grid-cols-[1.4fr_1fr]">
          <WalletCard />
          <Link
            to={ROUTES.logTransaction}
            className="group relative flex flex-col items-center justify-center gap-3 rounded-3xl border-2 border-dashed border-border bg-card/40 p-8 text-center transition-colors hover:border-coral hover:bg-coral/5"
          >
            <div className="grid h-14 w-14 place-items-center rounded-2xl bg-coral/15 text-coral transition-transform group-hover:scale-110">
              <Upload size={22} />
            </div>
            <div>
              <div className="text-lg">Drop a bill</div>
              <div className="mt-1 text-xs text-muted-foreground">PDF, JPG, or paste a screenshot</div>
            </div>
            <div className="mt-3 rounded-full bg-coral px-4 py-2 text-xs font-medium text-primary-foreground">
              Log transaction →
            </div>
          </Link>
        </div>

        <div className="mt-6 grid gap-6 lg:grid-cols-[1.4fr_1fr]">
          <RecentTransactions />
          <div className="flex flex-col gap-6">
            <PortfolioMiniChart />
          </div>
        </div>

        <div className="mt-6">
          <MarketSignalIndicator />
        </div>
      </div>
    </AppShell>
  );
}

