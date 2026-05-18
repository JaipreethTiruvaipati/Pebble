import { motion } from "framer-motion";
import { Info } from "lucide-react";
import { useState } from "react";
import { CategoryChip } from "@/components/ui/CategoryChip";
import { Slider } from "@/components/ui/Slider";
import { formatCurrency } from "@/lib/formatCurrency";
import { calcPenalty } from "@/lib/calcPenalty";
import type { DraftLineItem } from "@/stores/draftStore";
import { scoreColor } from "@/lib/scoreColor";

export function LineItemCard({ item, index }: { item: DraftLineItem; index: number }) {
  const [openTip, setOpenTip] = useState(false);
  const penalty = calcPenalty(item.amount, item.score);
  const color = scoreColor(item.score);

  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }}
      transition={{ delay: index * 0.08, duration: 0.5, ease: [0.16, 1, 0.3, 1] }}
      className="rounded-2xl border border-border bg-card p-5"
    >
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <h4 className="truncate text-base">{item.name}</h4>
            <CategoryChip category={item.category} />
          </div>
          <div className="mt-1 font-mono text-sm text-muted-foreground">{formatCurrency(item.amount)}</div>
        </div>
        <div className="text-right">
          <div className="font-mono text-2xl" style={{ color }}>{item.score}</div>
          <div className="text-[10px] uppercase tracking-widest text-muted-foreground">impulse</div>
        </div>
      </div>

      <div className="mt-4">
        <Slider value={item.score} height={6} />
      </div>

      <div className="mt-4 flex items-center justify-between">
        <button
          onMouseEnter={() => setOpenTip(true)}
          onMouseLeave={() => setOpenTip(false)}
          onFocus={() => setOpenTip(true)}
          onBlur={() => setOpenTip(false)}
          className="relative inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground"
        >
          <Info size={12} /> Why this score?
          {openTip && (
            <span className="absolute bottom-full left-0 mb-2 w-64 rounded-lg border border-border bg-card-elevated p-3 text-left text-[11px] leading-relaxed text-foreground shadow-2xl">
              {item.reasoning}
            </span>
          )}
        </button>
        <div className="font-mono text-sm">
          Penalty: <span className="text-coral">{formatCurrency(penalty)}</span>
        </div>
      </div>
    </motion.div>
  );
}
export default LineItemCard;
