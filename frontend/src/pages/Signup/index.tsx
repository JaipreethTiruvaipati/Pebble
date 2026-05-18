import { useState } from "react";
import { Link, useNavigate } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import { ROUTES } from "@/routes";
import { useAuthStore } from "@/stores/authStore";
import { toast } from "@/components/ui/Toast";
import { topupWallet } from "@/api/wallet.api";

interface StepProps {
  next: () => void;
  email: string;
  setEmail: (v: string) => void;
}

export function AccountStep({ next, email, setEmail }: StepProps) {
  return (
    <div className="space-y-4">
      <Input label="Full name" defaultValue="" placeholder="Your name" />
      <Input label="Email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
      <Input label="Mobile" type="tel" placeholder="+91 98123 45678" />
      <Button onClick={next} className="w-full" disabled={!email.includes("@")}>
        Continue
      </Button>
    </div>
  );
}

export function VerifyStep({ next }: Pick<StepProps, "next">) {
  return (
    <div className="space-y-4">
      <p className="text-sm text-muted-foreground">Dev mode: skip OTP — we use email login on the next step.</p>
      <Button onClick={next} className="w-full">
        Continue
      </Button>
    </div>
  );
}

export function FundStep({ email }: { email: string }) {
  const nav = useNavigate();
  const { login, isLoading } = useAuthStore();
  const [amount, setAmount] = useState("5000");

  const finish = async () => {
    try {
      await login(email, "signup");
      const topup = parseFloat(amount);
      if (topup > 0) await topupWallet(topup);
      toast("Account ready!", "success");
      nav({ to: ROUTES.onboardingRisk });
    } catch (err) {
      toast(err instanceof Error ? err.message : "Signup failed", "error");
    }
  };

  return (
    <div className="space-y-4">
      <Input label="UPI ID for top-ups" defaultValue="" placeholder="you@bank" />
      <Input
        label="Initial deposit (₹)"
        value={amount}
        onChange={(e) => setAmount(e.target.value)}
      />
      <Button onClick={() => void finish()} className="w-full" disabled={isLoading}>
        {isLoading ? "Setting up…" : "Finish & set up"}
      </Button>
    </div>
  );
}

export default function Signup() {
  const [step, setStep] = useState(0);
  const [email, setEmail] = useState("");
  const steps = [
    { title: "Create account", el: <AccountStep next={() => setStep(1)} email={email} setEmail={setEmail} /> },
    { title: "Verify mobile", el: <VerifyStep next={() => setStep(2)} /> },
    { title: "Fund wallet", el: <FundStep email={email} /> },
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
        <div className="mt-6 text-center text-xs text-muted-foreground">
          Already have an account?{" "}
          <Link to={ROUTES.login} className="text-coral">
            Log in
          </Link>
        </div>
      </motion.div>
    </div>
  );
}
