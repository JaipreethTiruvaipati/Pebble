// score 0 (best, smart spend) -> teal; 50 -> amber; 100 (worst, impulsive) -> coral
export function scoreColor(score: number): string {
  if (score <= 33) return "var(--teal)";
  if (score <= 66) return "var(--amber)";
  return "var(--coral)";
}

export function scoreLabel(score: number): string {
  if (score <= 33) return "Smart";
  if (score <= 66) return "Borderline";
  return "Impulsive";
}

export function scoreBgClass(score: number): string {
  if (score <= 33) return "bg-teal/15 text-teal";
  if (score <= 66) return "bg-amber/15 text-amber";
  return "bg-coral/15 text-coral";
}
