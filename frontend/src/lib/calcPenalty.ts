// Penalty = a fraction of amount, scaled non-linearly with impulse score.
// Score 0 -> 0 penalty. Score 100 -> 20% of amount.
export function calcPenalty(amount: number, score: number): number {
  const rate = Math.pow(Math.max(0, Math.min(100, score)) / 100, 1.5) * 0.2;
  return Math.round(amount * rate);
}
