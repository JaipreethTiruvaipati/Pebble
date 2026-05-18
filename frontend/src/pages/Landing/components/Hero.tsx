import { motion } from "framer-motion";
import { useEffect, useState } from "react";
import { Link } from "@tanstack/react-router";
import { Button } from "@/components/ui/Button";
import { ROUTES } from "@/routes";
import { ArrowRight } from "lucide-react";

const phrases = ["Spend like you mean it.", "Penalize the impulse.", "Compound the wisdom."];

export function Hero() {
  const [text, setText] = useState("");
  const [phraseIdx, setPhraseIdx] = useState(0);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    const cur = phrases[phraseIdx];
    const tick = setTimeout(() => {
      if (!deleting) {
        if (text.length < cur.length) setText(cur.slice(0, text.length + 1));
        else setTimeout(() => setDeleting(true), 1400);
      } else {
        if (text.length > 0) setText(cur.slice(0, text.length - 1));
        else { setDeleting(false); setPhraseIdx((p) => (p + 1) % phrases.length); }
      }
    }, deleting ? 40 : 70);
    return () => clearTimeout(tick);
  }, [text, deleting, phraseIdx]);

  return (
    <section className="relative isolate overflow-hidden">
      <div className="absolute inset-0 mesh-gradient" />
      <motion.div
        className="absolute -top-32 left-1/3 h-[40rem] w-[40rem] rounded-full bg-coral/30 blur-3xl"
        animate={{ x: [0, 60, -40, 0], y: [0, -30, 40, 0] }}
        transition={{ duration: 18, repeat: Infinity, ease: "easeInOut" }}
      />
      <motion.div
        className="absolute bottom-0 right-0 h-[36rem] w-[36rem] rounded-full bg-teal/20 blur-3xl"
        animate={{ x: [0, -50, 30, 0], y: [0, 30, -40, 0] }}
        transition={{ duration: 22, repeat: Infinity, ease: "easeInOut" }}
      />
      <div className="relative mx-auto flex min-h-screen max-w-6xl flex-col items-center justify-center px-6 text-center">
        <motion.span
          initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.6 }}
          className="mb-6 inline-flex items-center gap-2 rounded-full border border-border bg-card/60 px-4 py-1.5 text-xs font-medium uppercase tracking-widest text-muted-foreground backdrop-blur"
        >
          <span className="h-1.5 w-1.5 rounded-full bg-teal" /> Built for India's next gen investors
        </motion.span>
        <h1 className="font-display text-5xl md:text-7xl font-semibold leading-[1.05] tracking-tight max-w-4xl">
          The wallet that <span className="text-coral">fines</span> your impulses and <span className="text-teal">invests</span> the difference.
        </h1>
        <p className="mt-8 h-7 text-lg md:text-xl text-muted-foreground">
          <span>{text}</span>
          <span className="ml-1 inline-block h-5 w-0.5 translate-y-1 animate-pulse bg-coral" />
        </p>
        <div className="mt-10 flex flex-wrap items-center justify-center gap-3">
          <Button asChild size="lg" variant="coral"><Link to={ROUTES.signup}>Get started <ArrowRight size={16} /></Link></Button>
          <Button asChild size="lg" variant="outline"><Link to={ROUTES.dashboard}>See dashboard</Link></Button>
        </div>
        <motion.div
          initial={{ opacity: 0, y: 30 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.6, duration: 0.8 }}
          className="mt-20 grid grid-cols-3 gap-6 max-w-3xl w-full"
        >
          {[
            { k: "₹8,920", v: "auto-saved this month" },
            { k: "17", v: "impulses caught" },
            { k: "+12.4%", v: "portfolio return" },
          ].map((s) => (
            <div key={s.v} className="rounded-2xl border border-border bg-card/60 p-5 backdrop-blur text-left">
              <div className="font-mono text-2xl font-semibold">{s.k}</div>
              <div className="mt-1 text-xs text-muted-foreground">{s.v}</div>
            </div>
          ))}
        </motion.div>
      </div>
    </section>
  );
}
export default Hero;

