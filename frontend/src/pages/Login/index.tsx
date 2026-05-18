import { useState } from "react";
import { Link, useNavigate } from "@tanstack/react-router";
import { motion } from "framer-motion";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import { ROUTES } from "@/routes";
import { useAuthStore } from "@/stores/authStore";
import { toast } from "@/components/ui/Toast";

export default function Login() {
  const [email, setEmail] = useState("demo@pebble.in");
  const [password, setPassword] = useState("demo");
  const [referral, setReferral] = useState("");
  const { login, isLoading } = useAuthStore();
  const navigate = useNavigate();

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await login(email, password, referral || undefined);
      toast("Welcome back!", "success");
      navigate({ to: ROUTES.dashboard });
    } catch (err) {
      toast(err instanceof Error ? err.message : "Login failed", "error");
    }
  };

  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      className="relative grid min-h-screen place-items-center px-6"
    >
      <div className="absolute inset-0 mesh-gradient opacity-60" />
      <form
        onSubmit={onSubmit}
        className="relative w-full max-w-md rounded-3xl border border-border bg-card p-8"
      >
        <Link to={ROUTES.landing} className="font-display text-xl">
          Pebble
        </Link>
        <h2 className="mt-6 text-2xl">Log in</h2>
        <p className="mt-1 text-sm text-muted-foreground">
          Dev mode: any password works with your email.
        </p>
        <motion.div className="mt-6 space-y-4" style={{ display: "block" }}>
          <Input label="Email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} />
          <Input label="Password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
          <Input
            label="Referral code (optional)"
            value={referral}
            onChange={(e) => setReferral(e.target.value)}
          />
          <Button type="submit" className="w-full" disabled={isLoading}>
            {isLoading ? "Signing in…" : "Sign in"}
          </Button>
        </motion.div>
        <p className="mt-6 text-center text-xs text-muted-foreground">
          New here?{" "}
          <Link to={ROUTES.signup} className="text-coral">
            Create account
          </Link>
        </p>
      </form>
    </motion.div>
  );
}
