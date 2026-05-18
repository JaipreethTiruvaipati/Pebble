import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/Button";
import { useAuthStore, type RiskProfile } from "@/stores/authStore";
import { ROUTES } from "@/routes";

const profiles: { id: RiskProfile; title: string; body: string; allocation: string; accent: string }[] = [
  { id: "Conservative", title: "Conservative", body: "Capital protection first. Mostly bonds and gold.", allocation: "20/30/50", accent: "text-teal" },
  { id: "Moderate", title: "Moderate", body: "Balanced equity + debt. Steady compounding.", allocation: "55/20/25", accent: "text-amber" },
  { id: "Aggressive", title: "Aggressive", body: "Equity-heavy. Higher volatility, higher upside.", allocation: "80/10/10", accent: "text-coral" },
];

export default function RiskProfile() {
  const { setRiskProfile, user } = useAuthStore();
  const [selected, setSelected] = useState<RiskProfile>(user?.riskProfile ?? "Moderate");
  const nav = useNavigate();

  return (
    <div className="mx-auto max-w-3xl px-6 py-16">
      <div className="text-xs uppercase tracking-widest text-coral">Step 1 of 2</div>
      <h1 className="mt-2 text-4xl">Pick your risk profile</h1>
      <p className="mt-2 text-muted-foreground">This shapes where your penalties get invested.</p>
      <div className="mt-8 grid gap-4 md:grid-cols-3">
        {profiles.map((p) => (
          <motion.button
            key={p.id}
            whileHover={{ y: -4 }}
            onClick={() => setSelected(p.id)}
            className={`text-left rounded-2xl border p-6 transition-colors ${selected === p.id ? "border-coral bg-coral/10" : "border-border bg-card hover:bg-card-elevated"}`}
          >
            <div className={`text-sm uppercase tracking-widest ${p.accent}`}>{p.title}</div>
            <div className="mt-2 font-mono text-2xl">{p.allocation}</div>
            <div className="text-[10px] text-muted-foreground">equity / gold / debt</div>
            <p className="mt-3 text-sm text-muted-foreground">{p.body}</p>
          </motion.button>
        ))}
      </div>
      <div className="mt-8 flex justify-end">
        <Button onClick={() => { setRiskProfile(selected); nav({ to: ROUTES.onboardingWallet }); }}>Continue →</Button>
      </div>
    </div>
  );
}

