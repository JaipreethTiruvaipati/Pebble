import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { AlertCircle, Clock } from "lucide-react";
import { countdown } from "@/lib/formatDate";
import { formatCurrency } from "@/lib/formatCurrency";
import { usePenalties } from "@/hooks/usePenalty";

export function PendingPenaltyBanner() {
  const { data } = usePenalties("pending");
  const top = data?.[0];
  const [t, setT] = useState({ hours: 0, minutes: 0, seconds: 0 });

  useEffect(() => {
    if (!top?.expires_at) return;
    const tick = () => setT(countdown(top.expires_at));
    tick();
    const i = setInterval(tick, 1000);
    return () => clearInterval(i);
  }, [top?.expires_at]);

  if (!top) return null;

  return (
    <motion.div
      initial={{ opacity: 0, y: -8 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      className="flex flex-wrap items-center gap-4 rounded-2xl border border-coral/30 bg-coral/10 px-5 py-4"
    >
      <div className="grid h-10 w-10 place-items-center rounded-xl bg-coral/20 text-coral">
        <AlertCircle size={18} />
      </div>
      <div className="min-w-0 flex-1">
        <div className="text-sm">
          Pending penalty: <span className="font-mono text-coral">{formatCurrency(top.amount)}</span> from{" "}
          <span className="text-foreground">{top.item_name || top.merchant}</span>
        </div>
        <div className="mt-0.5 text-xs text-muted-foreground">Auto-invests when the consent timer ends.</div>
      </div>
      <div className="flex items-center gap-2 rounded-xl bg-background/40 px-3 py-2 font-mono text-sm">
        <Clock size={14} className="text-coral" />
        <span>
          {String(t.hours).padStart(2, "0")}:{String(t.minutes).padStart(2, "0")}:{String(t.seconds).padStart(2, "0")}
        </span>
      </div>
    </motion.div>
  );
}
