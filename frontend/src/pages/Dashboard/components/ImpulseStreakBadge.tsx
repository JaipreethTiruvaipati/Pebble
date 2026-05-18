import { motion } from "framer-motion";
import { Flame } from "lucide-react";
import { streak } from "@/lib/firebase";

export function ImpulseStreakBadge() {
  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }} animate={{ opacity: 1, scale: 1 }} transition={{ delay: 0.2, duration: 0.5 }}
      className="flex items-center gap-3 rounded-2xl border border-amber/30 bg-amber/10 px-4 py-3"
    >
      <div className="grid h-10 w-10 place-items-center rounded-xl bg-amber/20 text-amber">
        <Flame size={18} />
      </div>
      <div>
        <div className="font-mono text-lg">{streak.current}<span className="text-xs text-muted-foreground"> days</span></div>
        <div className="text-[11px] text-muted-foreground">smart streak · best {streak.best}</div>
      </div>
    </motion.div>
  );
}
export default ImpulseStreakBadge;

