import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { ROUTES } from "@/routes";
import { topupWallet } from "@/api/wallet.api";
import { toast } from "@/components/ui/Toast";

export default function WalletSetup() {
  const nav = useNavigate();
  const [amount, setAmount] = useState("25000");
  const [loading, setLoading] = useState(false);

  const finish = async () => {
    setLoading(true);
    try {
      const val = parseFloat(amount);
      if (val > 0) await topupWallet(val);
      toast("Wallet ready", "success");
      nav({ to: ROUTES.dashboard });
    } catch (err) {
      toast(err instanceof Error ? err.message : "Wallet setup failed", "error");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto max-w-md px-6 py-16">
      <div className="text-xs uppercase tracking-widest text-coral">Step 2 of 2</div>
      <h1 className="mt-2 text-4xl">Set up your wallet</h1>
      <p className="mt-2 text-muted-foreground">We'll auto-debit penalties from this account.</p>
      <div className="mt-8 space-y-4">
        <Input label="Linked bank account" defaultValue="HDFC ****4521" disabled />
        <Input
          label="Initial top-up (₹)"
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
        />
        <Button className="w-full" onClick={() => void finish()} disabled={loading}>
          {loading ? "Saving…" : "Enter Pebble →"}
        </Button>
      </div>
    </div>
  );
}
