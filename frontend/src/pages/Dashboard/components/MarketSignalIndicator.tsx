import { ScoreCircle } from "@/components/ui/ScoreCircle";
import { TrendingUp } from "lucide-react";
import { useMarketSignal } from "@/hooks/usePortfolio";
import { Skeleton } from "@/components/ui/Skeleton";

export function MarketSignalIndicator() {
  const { data, isLoading } = useMarketSignal();
  const score = Math.round(data?.composite_score ?? 0);
  const label = score >= 60 ? "Bullish" : score >= 40 ? "Neutral" : "Cautious";
  const note =
    data?.signals?.[0]?.indicator ??
    "Market signals refresh from Redis when market-poller is running.";

  if (isLoading) return <Skeleton className="h-48 rounded-3xl" />;

  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <div className="flex items-center justify-between">
        <div>
          <div className="text-xs uppercase tracking-widest text-muted-foreground">Market pulse</div>
          <div className="mt-1 text-lg">Opportunity score</div>
        </div>
        <span className="inline-flex items-center gap-1 rounded-full bg-teal/15 px-3 py-1 text-xs text-teal">
          <TrendingUp size={12} /> {label}
        </span>
      </div>
      <div className="mt-4 flex items-center gap-6">
        <ScoreCircle score={score} size={140} label="Pulse" sublabel="NIFTY 50" />
        <p className="flex-1 text-sm leading-relaxed text-muted-foreground">{note}</p>
      </div>
    </div>
  );
}
