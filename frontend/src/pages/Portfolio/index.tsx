import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { MetricCard } from "@/components/ui/MetricCard";
import { Button } from "@/components/ui/Button";
import { PieChart, Pie, Cell, ResponsiveContainer, LineChart, Line, XAxis, YAxis, Tooltip, CartesianGrid } from "recharts";
import { motion } from "framer-motion";
import { allocation, holdings, portfolioTimeline } from "@/lib/firebase";
import { formatCurrency } from "@/lib/formatCurrency";

function AllocationDonut() {
  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <div className="text-xs uppercase tracking-widest text-muted-foreground">Allocation</div>
      <div className="mt-4 flex items-center gap-6">
        <div className="h-48 w-48">
          <ResponsiveContainer>
            <PieChart>
              <Pie
                data={allocation}
                dataKey="value"
                innerRadius={56}
                outerRadius={86}
                paddingAngle={3}
                strokeWidth={0}
                animationDuration={1100}
              >
                {allocation.map((a) => <Cell key={a.name} fill={a.color} />)}
              </Pie>
            </PieChart>
          </ResponsiveContainer>
        </div>
        <ul className="flex-1 space-y-3">
          {allocation.map((a, i) => (
            <motion.li
              key={a.name}
              initial={{ opacity: 0, x: 8 }} animate={{ opacity: 1, x: 0 }} transition={{ delay: i * 0.08 }}
              className="flex items-center justify-between"
            >
              <span className="flex items-center gap-2 text-sm">
                <span className="h-2.5 w-2.5 rounded-full" style={{ background: a.color }} />
                {a.name}
              </span>
              <span className="font-mono text-sm">{a.value}%</span>
            </motion.li>
          ))}
        </ul>
      </div>
    </div>
  );
}

function HoldingsList() {
  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <h3 className="text-lg">Holdings</h3>
      <ul className="mt-4 divide-y divide-border">
        {holdings.map((h, i) => (
          <motion.li
            key={h.name}
            initial={{ opacity: 0, y: 6 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: i * 0.06 }}
            className="flex items-center justify-between py-3"
          >
            <div>
              <div className="text-sm">{h.name}</div>
              <div className="text-[11px] text-muted-foreground">{h.allocation}% of portfolio</div>
            </div>
            <div className="text-right">
              <div className="font-mono text-sm">{formatCurrency(h.value)}</div>
              <div className="font-mono text-[11px] text-teal">▲ {h.change}%</div>
            </div>
          </motion.li>
        ))}
      </ul>
    </div>
  );
}

function TimelineChart() {
  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <div className="flex items-end justify-between">
        <div>
          <div className="text-xs uppercase tracking-widest text-muted-foreground">Portfolio timeline</div>
          <div className="mt-2 font-mono text-2xl">{formatCurrency(86400)}</div>
        </div>
        <div className="text-teal font-mono text-sm">▲ 35.0% · 12W</div>
      </div>
      <div className="mt-6 h-64">
        <ResponsiveContainer>
          <LineChart data={portfolioTimeline}>
            <CartesianGrid stroke="#232745" strokeDasharray="3 3" vertical={false} />
            <XAxis dataKey="day" stroke="#8B8FA8" fontSize={11} axisLine={false} tickLine={false} />
            <YAxis stroke="#8B8FA8" fontSize={11} axisLine={false} tickLine={false} tickFormatter={(v) => `₹${(v / 1000).toFixed(0)}k`} />
            <Tooltip
              contentStyle={{ background: "#1C1F35", border: "1px solid #232745", borderRadius: 12, fontFamily: "JetBrains Mono" }}
              formatter={(v: any) => formatCurrency(Number(v))}
            />
            <Line dataKey="value" stroke="#E94560" strokeWidth={2.5} dot={false} type="monotone" animationDuration={1400} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}

export default function Portfolio() {
  return (
    <AppShell>
      <div className="mx-auto max-w-7xl p-6 md:p-10">
        <PageHeader
          title="Portfolio"
          subtitle="Where your penalties became investments."
          actions={<Button size="sm" variant="outline">Rebalance</Button>}
        />

        <div className="grid gap-4 md:grid-cols-3">
          <MetricCard label="Total invested" value={64000} currency accent="default" />
          <MetricCard label="Current value" value={86400} currency accent="teal" delta={{ value: 35, positive: true }} />
          <MetricCard label="Returns" value={22400} currency accent="teal" footer="₹3,840 from fines this month" />
        </div>

        <div className="mt-6 grid gap-6 lg:grid-cols-[1.4fr_1fr]">
          <TimelineChart />
          <AllocationDonut />
        </div>

        <div className="mt-6">
          <HoldingsList />
        </div>
      </div>
    </AppShell>
  );
}
