import { motion } from "framer-motion";
import { Flame } from "lucide-react";
import { streak } from "@/lib/firebase";

export function StreakCounter() {
  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.96 }} animate={{ opacity: 1, scale: 1 }}
      className="relative overflow-hidden rounded-3xl border border-border bg-card p-6"
    >
      <div className="absolute -right-8 -top-8 h-40 w-40 rounded-full bg-amber/10 blur-3xl" />
      <div className="relative">
        <div className="flex items-center gap-2 text-xs uppercase tracking-widest text-muted-foreground">
          <Flame size={14} className="text-amber" /> Smart streak
        </div>
        <div className="mt-3 flex items-end gap-3">
          <span className="font-mono text-6xl">{streak.current}</span>
          <span className="mb-2 text-sm text-muted-foreground">days</span>
        </div>
        <div className="mt-3 text-xs text-muted-foreground">Personal best: <span className="font-mono text-foreground">{streak.best} days</span></div>
        <div className="mt-5 flex gap-1">
          {Array.from({ length: 14 }).map((_, i) => (
            <motion.div
              key={i}
              initial={{ scaleY: 0 }} animate={{ scaleY: 1 }} transition={{ delay: i * 0.04 }}
              className={`h-8 flex-1 rounded ${i < streak.current ? "bg-amber" : "bg-card-elevated"}`}
            />
          ))}
        </div>
      </div>
    </motion.div>
  );
}
export default StreakCounter;

