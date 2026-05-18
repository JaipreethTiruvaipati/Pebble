import { motion } from "framer-motion";
import { peerBenchmark } from "@/lib/firebase";
import { scoreColor } from "@/lib/scoreColor";

export function PeerBenchmark() {
  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <h3 className="text-lg">Peer benchmark</h3>
      <p className="mt-1 text-xs text-muted-foreground">{peerBenchmark.cohort} cohort</p>
      <div className="mt-6 space-y-5">
        {[
          { label: "You", value: peerBenchmark.you },
          { label: "Peers", value: peerBenchmark.peers },
        ].map((row, i) => (
          <div key={row.label}>
            <div className="flex items-center justify-between text-sm">
              <span>{row.label}</span>
              <span className="font-mono" style={{ color: scoreColor(row.value) }}>{row.value}</span>
            </div>
            <div className="mt-1.5 h-3 overflow-hidden rounded-full bg-card-elevated">
              <motion.div
                initial={{ width: 0 }} animate={{ width: `${row.value}%` }}
                transition={{ delay: i * 0.15, duration: 1, ease: [0.16, 1, 0.3, 1] }}
                className="h-full rounded-full"
                style={{ background: scoreColor(row.value) }}
              />
            </div>
          </div>
        ))}
      </div>
      <div className="mt-5 rounded-xl bg-teal/10 p-3 text-xs text-teal">
        You're {peerBenchmark.peers - peerBenchmark.you} points more disciplined than peers.
      </div>
    </div>
  );
}
export default PeerBenchmark;

