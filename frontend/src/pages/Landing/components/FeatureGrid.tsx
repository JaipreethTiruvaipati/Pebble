import { motion } from "framer-motion";
import { ScanLine, Brain, TrendingUp } from "lucide-react";

const features = [
  {
    icon: ScanLine,
    title: "Snap any bill",
    body: "Upload a receipt and Pebble extracts every line item, instantly. No tagging, no spreadsheets.",
    accent: "text-coral",
  },
  {
    icon: Brain,
    title: "Impulse scoring",
    body: "Each item is scored 0–100 for impulsiveness. Smart spends pass through; impulses get fined.",
    accent: "text-amber",
  },
  {
    icon: TrendingUp,
    title: "Auto-invest the fines",
    body: "Every penalty routes straight into a curated portfolio matched to your risk profile.",
    accent: "text-teal",
  },
];

export function FeatureGrid() {
  return (
    <section className="relative mx-auto max-w-6xl px-6 py-24">
      <div className="max-w-2xl">
        <h2 className="text-4xl font-semibold tracking-tight">A wallet that argues with you.</h2>
        <p className="mt-3 text-muted-foreground">Three things working in the background, every time you spend.</p>
      </div>
      <div className="mt-12 grid gap-5 md:grid-cols-3">
        {features.map((f, i) => (
          <motion.div
            key={f.title}
            initial={{ opacity: 0, y: 24 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true, amount: 0.3 }}
            transition={{ delay: i * 0.1, duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
            className="group relative overflow-hidden rounded-3xl border border-border bg-card p-7"
          >
            <div className={`mb-5 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-card-elevated ${f.accent}`}>
              <f.icon size={22} />
            </div>
            <h3 className="text-xl">{f.title}</h3>
            <p className="mt-3 text-sm leading-relaxed text-muted-foreground">{f.body}</p>
            <div className="pointer-events-none absolute -bottom-12 -right-12 h-40 w-40 rounded-full bg-coral/10 blur-3xl opacity-0 transition-opacity group-hover:opacity-100" />
          </motion.div>
        ))}
      </div>
    </section>
  );
}
export default FeatureGrid;
