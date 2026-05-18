import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { AlertCircle, Clock } from "lucide-react";
import { pendingPenalty } from "@/lib/firebase";
import { countdown } from "@/lib/formatDate";
import { formatCurrency } from "@/lib/formatCurrency";

export function PendingPenaltyBanner() {
  const [t, setT] = useState(countdown(pendingPenalty.expiresAt));
  useEffect(() => {
    const i = setInterval(() => setT(countdown(pendingPenalty.expiresAt)), 1000);
    return () => clearInterval(i);
  }, []);

  return (
    <motion.div
      initial={{ opacity: 0, y: -8 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.5 }}
      className="flex flex-wrap items-center gap-4 rounded-2xl border border-coral/30 bg-coral/10 px-5 py-4"
    >
      <div className="grid h-10 w-10 place-items-center rounded-xl bg-coral/20 text-coral">
        <AlertCircle size={18} />
      </div>
      <div className="flex-1 min-w-0">
        <div className="text-sm">
          Pending penalty: <span className="font-mono text-coral">{formatCurrency(pendingPenalty.amount)}</span> from <span className="text-foreground">{pendingPenalty.source}</span>
        </div>
        <div className="mt-0.5 text-xs text-muted-foreground">Auto-routes to investments when timer ends.</div>
      </div>
      <div className="flex items-center gap-2 rounded-xl bg-background/40 px-3 py-2 font-mono text-sm">
        <Clock size={14} className="text-coral" />
        <span>{String(t.hours).padStart(2, "0")}:{String(t.minutes).padStart(2, "0")}:{String(t.seconds).padStart(2, "0")}</span>
      </div>
      <button className="rounded-lg border border-coral/40 px-3 py-2 text-xs hover:bg-coral/20">Review</button>
    </motion.div>
  );
}
export default PendingPenaltyBanner;

