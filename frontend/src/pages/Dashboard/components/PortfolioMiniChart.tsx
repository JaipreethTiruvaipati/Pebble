import { AreaChart, Area, ResponsiveContainer, Tooltip } from "recharts";
import { portfolioTimeline } from "@/lib/firebase";
import { formatCurrency } from "@/lib/formatCurrency";

export function PortfolioMiniChart() {
  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <div className="flex items-end justify-between">
        <div>
          <div className="text-xs uppercase tracking-widest text-muted-foreground">Portfolio value</div>
          <div className="mt-2 font-mono text-3xl">{formatCurrency(86400)}</div>
        </div>
        <div className="text-right">
          <div className="text-teal font-mono text-sm">▲ ₹2,200</div>
          <div className="text-[11px] text-muted-foreground">last 12 weeks</div>
        </div>
      </div>
      <div className="mt-4 h-40">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={portfolioTimeline}>
            <defs>
              <linearGradient id="mini" x1="0" x2="0" y1="0" y2="1">
                <stop offset="0%" stopColor="#00D4AA" stopOpacity={0.4} />
                <stop offset="100%" stopColor="#00D4AA" stopOpacity={0} />
              </linearGradient>
            </defs>
            <Tooltip
              cursor={{ stroke: "#232745" }}
              contentStyle={{ background: "#1C1F35", border: "1px solid #232745", borderRadius: 12, fontFamily: "JetBrains Mono" }}
              labelStyle={{ color: "#8B8FA8" }}
              formatter={(v: any) => formatCurrency(Number(v))}
            />
            <Area dataKey="value" stroke="#00D4AA" strokeWidth={2} fill="url(#mini)" type="monotone" />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
export default PortfolioMiniChart;

