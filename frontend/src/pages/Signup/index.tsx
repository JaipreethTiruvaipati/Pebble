import { useState } from "react";
import { Link, useNavigate } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import { ROUTES } from "@/routes";

interface StepProps { next: () => void; }

export function AccountStep({ next }: StepProps) {
  return (
    <div className="space-y-4">
      <Input label="Full name" defaultValue="Arjun Sharma" />
      <Input label="Email" type="email" defaultValue="arjun.sharma@pebble.in" />
      <Input label="Mobile" type="tel" defaultValue="+91 98123 45678" />
      <Button onClick={next} className="w-full">Continue</Button>
    </div>
  );
}

export function VerifyStep({ next }: StepProps) {
  return (
    <div className="space-y-4">
      <p className="text-sm text-muted-foreground">We sent an OTP to your mobile. Enter it below.</p>
      <div className="flex gap-2">
        {[0,1,2,3,4,5].map((i) => (
          <input key={i} maxLength={1} defaultValue={"482917"[i]} className="h-14 w-12 rounded-xl bg-input border border-border text-center font-mono text-xl outline-none focus:border-coral" />
        ))}
      </div>
      <Button onClick={next} className="w-full">Verify &amp; continue</Button>
    </div>
  );
}

export function FundStep() {
  const nav = useNavigate();
  return (
    <div className="space-y-4">
      <Input label="UPI ID for top-ups" defaultValue="arjun@hdfc" />
      <Input label="Initial deposit (₹)" defaultValue="5000" />
      <Button onClick={() => nav({ to: ROUTES.onboardingRisk })} className="w-full">Finish &amp; set up</Button>
    </div>
  );
}

export default function Signup() {
  const [step, setStep] = useState(0);
  const steps = [
    { title: "Create account", el: <AccountStep next={() => setStep(1)} /> },
    { title: "Verify mobile", el: <VerifyStep next={() => setStep(2)} /> },
    { title: "Fund wallet", el: <FundStep /> },
  ];
  return (
    <div className="relative grid min-h-screen place-items-center px-6">
      <div className="absolute inset-0 mesh-gradient opacity-60" />
      <motion.div initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }} className="relative w-full max-w-md rounded-3xl border border-border bg-card p-8">
        <Link to={ROUTES.landing} className="font-display text-xl">Pebble</Link>
        <div className="mt-6">
          <div className="flex gap-1.5">
            {steps.map((_, i) => (
              <div key={i} className={`h-1 flex-1 rounded-full ${i <= step ? "bg-coral" : "bg-card-elevated"}`} />
            ))}
          </div>
          <h2 className="mt-6 text-2xl">{steps[step].title}</h2>
          <div className="mt-6">{steps[step].el}</div>
        </div>
        <div className="mt-6 text-center text-xs text-muted-foreground">Already have an account? <Link to={ROUTES.dashboard} className="text-coral">Log in</Link></div>
      </motion.div>
    </div>
  );
}

