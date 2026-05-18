import { useState } from "react";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Input } from "@/components/ui/Input";
import { CategoryChip } from "@/components/ui/CategoryChip";
import { recentTransactions } from "@/lib/firebase";
import { formatCurrency } from "@/lib/formatCurrency";
import { formatDate } from "@/lib/formatDate";
import { scoreBgClass } from "@/lib/scoreColor";

export default function TransactionHistory() {
  const [q, setQ] = useState("");
  const filtered = recentTransactions.filter((t) => t.merchant.toLowerCase().includes(q.toLowerCase()));

  return (
    <AppShell>
      <div className="mx-auto max-w-5xl p-6 md:p-10">
        <PageHeader title="Transaction history" subtitle="Every spend, scored." />
        <div className="mb-4 max-w-sm">
          <Input placeholder="Search merchant…" value={q} onChange={(e) => setQ(e.target.value)} />
        </div>
        <div className="overflow-hidden rounded-3xl border border-border bg-card">
          <table className="w-full text-sm">
            <thead className="bg-card-elevated text-left text-xs uppercase tracking-widest text-muted-foreground">
              <tr>
                <th className="p-4">Date</th><th className="p-4">Merchant</th><th className="p-4">Category</th>
                <th className="p-4 text-right">Amount</th><th className="p-4 text-right">Penalty</th><th className="p-4 text-right">Score</th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((t) => (
                <tr key={t.id} className="border-t border-border hover:bg-card-elevated/50">
                  <td className="p-4 text-muted-foreground">{formatDate(t.date)}</td>
                  <td className="p-4">{t.merchant}</td>
                  <td className="p-4"><CategoryChip category={t.category} /></td>
                  <td className="p-4 text-right font-mono">{formatCurrency(t.amount)}</td>
                  <td className="p-4 text-right font-mono text-coral">{t.penalty ? `−${formatCurrency(t.penalty)}` : "—"}</td>
                  <td className="p-4 text-right"><span className={`inline-block rounded-full px-2 py-0.5 font-mono text-[11px] ${scoreBgClass(t.score)}`}>{t.score}</span></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </AppShell>
  );
}

