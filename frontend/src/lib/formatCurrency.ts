export function formatCurrency(value: number, opts: { showDecimals?: boolean; compact?: boolean } = {}): string {
  const { showDecimals = false, compact = false } = opts;
  if (compact && Math.abs(value) >= 100000) {
    if (Math.abs(value) >= 10000000) return `₹${(value / 10000000).toFixed(2)}Cr`;
    return `₹${(value / 100000).toFixed(2)}L`;
  }
  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: "INR",
    minimumFractionDigits: showDecimals ? 2 : 0,
    maximumFractionDigits: showDecimals ? 2 : 0,
  }).format(value);
}

export function formatNumber(value: number): string {
  return new Intl.NumberFormat("en-IN").format(value);
}
