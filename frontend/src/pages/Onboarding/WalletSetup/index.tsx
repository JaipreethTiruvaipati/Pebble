import { useNavigate } from "@tanstack/react-router";
import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { ROUTES } from "@/routes";

export default function WalletSetup() {
  const nav = useNavigate();
  return (
    <div className="mx-auto max-w-md px-6 py-16">
      <div className="text-xs uppercase tracking-widest text-coral">Step 2 of 2</div>
      <h1 className="mt-2 text-4xl">Set up your wallet</h1>
      <p className="mt-2 text-muted-foreground">We'll auto-debit penalties from this account.</p>
      <div className="mt-8 space-y-4">
        <Input label="Linked bank account" defaultValue="HDFC ****4521" />
        <Input label="Monthly invest target (₹)" defaultValue="25000" />
        <Input label="Auto-debit consent" defaultValue="UPI Mandate · Active" disabled />
        <Button className="w-full" onClick={() => nav({ to: ROUTES.dashboard })}>Enter Pebble →</Button>
      </div>
    </div>
  );
}

