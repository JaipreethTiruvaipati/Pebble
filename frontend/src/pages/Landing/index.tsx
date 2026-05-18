import { Link } from "@tanstack/react-router";
import { Hero } from "./components/Hero";
import { FeatureGrid } from "./components/FeatureGrid";
import { CTASection } from "./components/CTASection";
import { ROUTES } from "@/routes";
import { motion } from "framer-motion";
import { AreaChart, Area, ResponsiveContainer } from "recharts";
import { portfolioTimeline } from "@/lib/firebase";

const steps = [
  { n: "01", title: "Snap", body: "Upload any bill or receipt." },
  { n: "02", title: "Score", body: "AI rates each item's impulsiveness." },
  { n: "03", title: "Invest", body: "Penalties auto-route into your portfolio." },
];

function StepFlow() {
  return (
    <section className="mx-auto max-w-6xl px-6 py-24">
      <h2 className="max-w-xl text-4xl font-semibold tracking-tight">From swipe to SIP, in three steps.</h2>
      <div className="mt-12 grid gap-5 md:grid-cols-3">
        {steps.map((s, i) => (
          <motion.div
            key={s.n}
            initial={{ opacity: 0, y: 20 }} whileInView={{ opacity: 1, y: 0 }} viewport={{ once: true }}
            transition={{ delay: i * 0.12, duration: 0.6 }}
            className="relative rounded-2xl border border-border bg-card p-7"
          >
            <div className="font-mono text-coral text-sm">{s.n}</div>
            <div className="mt-3 text-2xl font-display">{s.title}</div>
            <div className="mt-2 text-sm text-muted-foreground">{s.body}</div>
          </motion.div>
        ))}
      </div>
    </section>
  );
}

function GrowthChart() {
  return (
    <section className="mx-auto max-w-6xl px-6 py-12">
      <div className="rounded-3xl border border-border bg-card p-7">
        <div className="flex items-end justify-between">
          <div>
            <div className="text-xs uppercase tracking-widest text-muted-foreground">Sample portfolio growth · 12 weeks</div>
            <div className="mt-2 font-mono text-3xl">₹86,400</div>
          </div>
          <div className="text-teal font-mono text-sm">▲ 35.0%</div>
        </div>
        <div className="mt-6 h-48">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={portfolioTimeline}>
              <defs>
                <linearGradient id="land" x1="0" x2="0" y1="0" y2="1">
                  <stop offset="0%" stopColor="#E94560" stopOpacity={0.5} />
                  <stop offset="100%" stopColor="#E94560" stopOpacity={0} />
                </linearGradient>
              </defs>
              <Area dataKey="value" stroke="#E94560" strokeWidth={2.5} fill="url(#land)" type="monotone" />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </div>
    </section>
  );
}

function TopBar() {
  return (
    <div className="absolute top-0 left-0 right-0 z-20 flex items-center justify-between px-6 py-5 max-w-6xl mx-auto">
      <Link to={ROUTES.landing} className="font-display text-xl tracking-tight">Pebble</Link>
      <nav className="flex items-center gap-6 text-sm text-muted-foreground">
        <Link to={ROUTES.dashboard} className="hover:text-foreground">Dashboard</Link>
        <Link to={ROUTES.signup} className="rounded-full bg-coral px-4 py-1.5 text-primary-foreground hover:brightness-110">Sign up</Link>
      </nav>
    </div>
  );
}

export default function Landing() {
  return (
    <div className="relative">
      <TopBar />
      <Hero />
      <StepFlow />
      <GrowthChart />
      <FeatureGrid />
      <CTASection />
      <footer className="border-t border-border py-8 text-center text-xs text-muted-foreground">
        © 2026 Pebble · Built with 🪨 in Bengaluru
      </footer>
    </div>
  );
}


