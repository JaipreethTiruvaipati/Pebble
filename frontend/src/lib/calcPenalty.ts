/** Penalty aligned with backend formula (min ₹5, max ₹500). */
export function calcPenalty(
  amount: number,
  score: number,
  penaltyRate = 0.1,
  threshold = 50,
): number {
  if (score < threshold) return 0;
  let raw = amount * penaltyRate * (Math.max(0, Math.min(100, score)) / 100);
  if (raw < 5) return 0;
  if (raw > 500) raw = 500;
  return Math.round(raw * 100) / 100;
}
