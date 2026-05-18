import { motion } from "framer-motion";
import { Wallet, ArrowUpRight } from "lucide-react";
import { CountUp } from "@/components/ui/MetricCard";
import { useWallet } from "@/hooks/useWallet";
import { Skeleton } from "@/components/ui/Skeleton";
import { toast } from "@/components/ui/Toast";

export function WalletCard() {
  const { balance, topup } = useWallet();
  const data = balance.data;
  const monthlyTarget = 25000;
  const invested = data?.invested_total ?? 0;
  const pct = monthlyTarget > 0 ? Math.min(100, Math.round((invested / monthlyTarget) * 100)) : 0;

  const onTopup = async () => {
    try {
      await topup.mutateAsync(5000);
      toast("Wallet topped up ₹5,000", "success");
    } catch {
      toast("Top-up failed", "error");
    }
  };

  if (balance.isLoading) return <Skeleton className="h-64 rounded-3xl" />;

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
      className="relative overflow-hidden rounded-3xl border border-border bg-card p-7"
    >
      <motion.div className="absolute -top-32 -right-24 h-72 w-72 rounded-full bg-coral/20 blur-3xl" />
      <div className="relative">
        <div className="flex items-center gap-2 text-xs uppercase tracking-widest text-muted-foreground">
          <Wallet size={14} /> Pebble wallet
        </div>
        <div className="mt-4 flex items-end gap-3">
          <CountUp value={data?.balance ?? 0} currency className="text-5xl md:text-6xl font-semibold tracking-tight" />
        </div>
        <div className="mt-8 grid grid-cols-3 gap-4">
          <div>
            <div className="text-[10px] uppercase tracking-widest text-muted-foreground">Invested</div>
            <CountUp value={invested} currency className="mt-1 block text-lg" />
          </div>
          <div>
            <div className="text-[10px] uppercase tracking-widest text-muted-foreground">Pending</div>
            <CountUp value={data?.pending_total ?? 0} currency className="mt-1 block text-lg text-coral" />
          </div>
          <div>
            <div className="text-[10px] uppercase tracking-widest text-muted-foreground">Total invested</div>
            <CountUp value={data?.invested_total ?? 0} currency className="mt-1 block text-lg text-teal" />
          </div>
        </div>
        <div className="mt-6">
          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <span>Invest progress</span>
            <span className="font-mono">{pct}%</span>
          </div>
          <div className="mt-2 h-2 overflow-hidden rounded-full bg-card-elevated">
            <motion.div
              initial={{ width: 0 }}
              animate={{ width: `${pct}%` }}
              transition={{ duration: 1.2, ease: [0.16, 1, 0.3, 1] }}
              className="h-full rounded-full bg-gradient-to-r from-coral to-amber"
            />
          </div>
        </div>
        <button
          type="button"
          onClick={onTopup}
          disabled={topup.isPending}
          className="mt-8 inline-flex items-center gap-2 rounded-xl bg-coral px-5 py-3 text-sm font-medium shadow-lg shadow-coral/30 hover:brightness-110 disabled:opacity-60"
        >
          Top up wallet <ArrowUpRight size={14} />
        </button>
      </div>
    </motion.div>
  );
}
