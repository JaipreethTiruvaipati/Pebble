import { motion } from "framer-motion";
import { scoreColor } from "@/lib/scoreColor";
import { useBenchmark } from "@/hooks/usePortfolio";
import { Skeleton } from "@/components/ui/Skeleton";

export function PeerBenchmark() {
  const { data, isLoading } = useBenchmark();
  if (isLoading || !data) return <Skeleton className="h-56 rounded-3xl" />;

  const you = Math.round(data.user_impulse_pct);
  const peers = Math.round(data.cohort_impulse_pct);

  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <h3 className="text-lg">Peer benchmark</h3>
      <p className="mt-1 text-xs text-muted-foreground">{data.cohort_label}</p>
      <div className="mt-6 space-y-5">
        {[
          { label: "You", value: you },
          { label: "Cohort", value: peers },
        ].map((row, i) => (
          <div key={row.label}>
            <div className="flex items-center justify-between text-sm">
              <span>{row.label}</span>
              <span className="font-mono" style={{ color: scoreColor(row.value) }}>
                {row.value}
              </span>
            </div>
            <div className="mt-1.5 h-3 overflow-hidden rounded-full bg-card-elevated">
              <motion.div
                initial={{ width: 0 }}
                animate={{ width: `${row.value}%` }}
                transition={{ delay: i * 0.15, duration: 1, ease: [0.16, 1, 0.3, 1] }}
                className="h-full rounded-full"
                style={{ background: scoreColor(row.value) }}
              />
            </div>
          </div>
        ))}
      </div>
      <div className="mt-5 rounded-xl bg-teal/10 p-3 text-xs text-teal">
        {data.saved_vs_cohort_pct > 0
          ? `You saved ${data.saved_vs_cohort_pct.toFixed(0)}% vs cohort impulse rate.`
          : "Keep logging transactions to improve your benchmark."}
      </div>
    </div>
  );
}
