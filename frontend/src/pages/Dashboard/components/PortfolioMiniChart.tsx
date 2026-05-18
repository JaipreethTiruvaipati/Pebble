import { AreaChart, Area, ResponsiveContainer, Tooltip } from "recharts";
import { formatCurrency } from "@/lib/formatCurrency";
import { usePortfolio } from "@/hooks/usePortfolio";
import { Skeleton } from "@/components/ui/Skeleton";

export function PortfolioMiniChart() {
  const { data, isLoading } = usePortfolio();
  const total = (data?.equity_value ?? 0) + (data?.gold_value ?? 0) + (data?.bond_value ?? 0);
  const chartData = [
    { day: "Eq", value: data?.equity_value ?? 0 },
    { day: "Au", value: data?.gold_value ?? 0 },
    { day: "Bd", value: data?.bond_value ?? 0 },
  ];

  if (isLoading) return <Skeleton className="h-56 rounded-3xl" />;

  return (
    <div className="rounded-3xl border border-border bg-card p-6">
      <div className="flex items-end justify-between">
        <div>
          <div className="text-xs uppercase tracking-widest text-muted-foreground">Portfolio value</div>
          <div className="mt-2 font-mono text-3xl">{formatCurrency(total || (data?.total_invested ?? 0))}</div>
        </div>
        <div className="text-right">
          <div className="font-mono text-sm text-teal">▲ {data?.gain_pct?.toFixed(1) ?? 0}%</div>
          <div className="text-[11px] text-muted-foreground">allocation mix</div>
        </div>
      </div>
      <div className="mt-4 h-40">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={chartData}>
            <defs>
              <linearGradient id="mini" x1="0" x2="0" y1="0" y2="1">
                <stop offset="0%" stopColor="#00D4AA" stopOpacity={0.4} />
                <stop offset="100%" stopColor="#00D4AA" stopOpacity={0} />
              </linearGradient>
            </defs>
            <Tooltip
              cursor={{ stroke: "#232745" }}
              contentStyle={{
                background: "#1C1F35",
                border: "1px solid #232745",
                borderRadius: 12,
                fontFamily: "JetBrains Mono",
              }}
              formatter={(v: number) => formatCurrency(Number(v))}
            />
            <Area dataKey="value" stroke="#00D4AA" strokeWidth={2} fill="url(#mini)" type="monotone" />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
