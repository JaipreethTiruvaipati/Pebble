import { motion } from "framer-motion";
import { Wallet, ArrowUpRight } from "lucide-react";
import { CountUp } from "@/components/ui/MetricCard";
import { wallet } from "@/lib/firebase";

export function WalletCard() {
  const pct = Math.round((wallet.investedThisMonth / wallet.monthlyTarget) * 100);
  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.6, ease: [0.16, 1, 0.3, 1] }}
      className="relative overflow-hidden rounded-3xl border border-border bg-card p-7"
    >
      <div className="absolute -top-32 -right-24 h-72 w-72 rounded-full bg-coral/20 blur-3xl" />
      <div className="absolute -bottom-32 -left-24 h-64 w-64 rounded-full bg-teal/10 blur-3xl" />
      <div className="relative">
        <div className="flex items-center gap-2 text-xs uppercase tracking-widest text-muted-foreground">
          <Wallet size={14} /> Pebble wallet
        </div>
        <div className="mt-4 flex items-end gap-3">
          <CountUp value={wallet.balance} currency className="text-5xl md:text-6xl font-semibold tracking-tight" />
          <span className="mb-2 text-xs text-teal font-mono">▲ ₹2,840 today</span>
        </div>

        <div className="mt-8 grid grid-cols-3 gap-4">
          <div>
            <div className="text-[10px] uppercase tracking-widest text-muted-foreground">Invested · May</div>
            <CountUp value={wallet.investedThisMonth} currency className="mt-1 block text-lg" />
          </div>
          <div>
            <div className="text-[10px] uppercase tracking-widest text-muted-foreground">Impulse saves</div>
            <CountUp value={wallet.impulseSavesCount} className="mt-1 block text-lg" />
          </div>
          <div>
            <div className="text-[10px] uppercase tracking-widest text-muted-foreground">Saved this month</div>
            <CountUp value={wallet.impulseSavedAmount} currency className="mt-1 block text-lg text-teal" />
          </div>
        </div>

        <div className="mt-6">
          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <span>Monthly invest target</span>
            <span className="font-mono">₹{wallet.investedThisMonth.toLocaleString("en-IN")} / ₹{wallet.monthlyTarget.toLocaleString("en-IN")}</span>
          </div>
          <div className="mt-2 h-2 overflow-hidden rounded-full bg-card-elevated">
            <motion.div
              initial={{ width: 0 }} animate={{ width: `${pct}%` }} transition={{ duration: 1.2, ease: [0.16, 1, 0.3, 1] }}
              className="h-full rounded-full bg-gradient-to-r from-coral to-amber"
            />
          </div>
        </div>

        <button className="mt-8 inline-flex items-center gap-2 rounded-xl bg-coral px-5 py-3 text-sm font-medium shadow-lg shadow-coral/30 hover:brightness-110">
          Top up wallet <ArrowUpRight size={14} />
        </button>
      </div>
    </motion.div>
  );
}
export default WalletCard;

