import { ScoreCircle } from "@/components/ui/ScoreCircle";
import { marketPulse } from "@/lib/firebase";
import { TrendingUp } from "lucide-react";

export function MarketSignalIndicator() {
  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <div className="flex items-center justify-between">
        <div>
          <div className="text-xs uppercase tracking-widest text-muted-foreground">Market pulse</div>
          <div className="mt-1 text-lg">Opportunity score</div>
        </div>
        <span className="inline-flex items-center gap-1 rounded-full bg-teal/15 px-3 py-1 text-xs text-teal">
          <TrendingUp size={12} /> {marketPulse.label}
        </span>
      </div>
      <div className="mt-4 flex items-center gap-6">
        <ScoreCircle score={marketPulse.score} size={140} label="Pulse" sublabel="NIFTY 50" />
        <p className="flex-1 text-sm text-muted-foreground leading-relaxed">{marketPulse.note}</p>
      </div>
    </div>
  );
}
export default MarketSignalIndicator;

