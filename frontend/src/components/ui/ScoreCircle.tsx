import { motion } from "framer-motion";
import { scoreColor, scoreLabel } from "@/lib/scoreColor";

export function ScoreCircle({ score, size = 160, label, sublabel }: { score: number; size?: number; label?: string; sublabel?: string }) {
  const radius = (size - 16) / 2;
  const circ = 2 * Math.PI * radius;
  const offset = circ - (score / 100) * circ;
  const color = scoreColor(score);
  return (
    <div className="relative inline-flex flex-col items-center" style={{ width: size, height: size }}>
      <svg width={size} height={size} className="-rotate-90">
        <circle cx={size / 2} cy={size / 2} r={radius} stroke="var(--card-elevated)" strokeWidth={10} fill="none" />
        <motion.circle
          cx={size / 2} cy={size / 2} r={radius}
          stroke={color} strokeWidth={10} strokeLinecap="round" fill="none"
          strokeDasharray={circ}
          initial={{ strokeDashoffset: circ }}
          animate={{ strokeDashoffset: offset }}
          transition={{ duration: 1.4, ease: [0.16, 1, 0.3, 1] }}
        />
      </svg>
      <div className="absolute inset-0 flex flex-col items-center justify-center">
        <motion.span
          className="font-mono text-4xl font-semibold"
          style={{ color }}
          initial={{ opacity: 0, y: 6 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.3, duration: 0.5 }}
        >
          {score}
        </motion.span>
        <span className="text-[11px] uppercase tracking-widest text-muted-foreground mt-0.5">{label ?? scoreLabel(score)}</span>
        {sublabel && <span className="text-[10px] text-muted-foreground mt-1">{sublabel}</span>}
      </div>
    </div>
  );
}
export default ScoreCircle;
