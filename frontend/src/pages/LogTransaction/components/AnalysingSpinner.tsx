import { motion } from "framer-motion";
import { Check, ScanText, Sparkles, Loader2 } from "lucide-react";
import { useEffect, useState } from "react";

const steps = [
  { label: "Reading bill", icon: ScanText },
  { label: "Extracting items", icon: Loader2 },
  { label: "Scoring impulses", icon: Sparkles },
];

export function AnalysingSpinner({ onDone }: { onDone: () => void }) {
  const [idx, setIdx] = useState(0);
  useEffect(() => {
    if (idx >= steps.length) { const t = setTimeout(onDone, 400); return () => clearTimeout(t); }
    const t = setTimeout(() => setIdx((i) => i + 1), 900);
    return () => clearTimeout(t);
  }, [idx, onDone]);

  return (
    <div className="mx-auto max-w-md rounded-3xl border border-border bg-card p-10">
      <h3 className="text-center text-xl">Analysing your bill</h3>
      <ul className="mt-8 space-y-4">
        {steps.map((s, i) => {
          const done = i < idx;
          const active = i === idx;
          return (
            <li key={s.label} className="flex items-center gap-4">
              <motion.div
                animate={{ scale: active ? 1.05 : 1 }}
                className={`grid h-10 w-10 place-items-center rounded-xl ${done ? "bg-teal/20 text-teal" : active ? "bg-coral/20 text-coral" : "bg-card-elevated text-muted-foreground"}`}
              >
                {done ? <Check size={16} /> : active ? <s.icon size={16} className="animate-spin" /> : <s.icon size={16} />}
              </motion.div>
              <span className={`text-sm ${done ? "text-teal" : active ? "text-foreground" : "text-muted-foreground"}`}>{s.label}</span>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
export default AnalysingSpinner;
