import { motion } from "framer-motion";
import { scoreColor } from "@/lib/scoreColor";

export function Slider({ value, max = 100, animate = true, height = 8, showLabel = false, className = "" }: { value: number; max?: number; animate?: boolean; height?: number; showLabel?: boolean; className?: string }) {
  const pct = Math.max(0, Math.min(100, (value / max) * 100));
  const color = scoreColor(value);
  return (
    <div className={className}>
      <div className="relative w-full overflow-hidden rounded-full bg-card-elevated" style={{ height }}>
        <motion.div
          className="absolute left-0 top-0 h-full rounded-full"
          style={{ backgroundColor: color }}
          initial={animate ? { width: 0 } : { width: `${pct}%` }}
          animate={{ width: `${pct}%` }}
          transition={{ duration: 1.1, ease: [0.16, 1, 0.3, 1] }}
        />
      </div>
      {showLabel && <div className="mt-1 text-right font-mono text-xs" style={{ color }}>{value}</div>}
    </div>
  );
}
export default Slider;
