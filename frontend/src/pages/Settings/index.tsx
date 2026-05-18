import { useEffect, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Toggle } from "@/components/ui/Toggle";
import { Input } from "@/components/ui/Input";
import { Button } from "@/components/ui/Button";
import { useAuthStore } from "@/stores/authStore";
import { getReferralMe } from "@/api/referral.api";
import { formatCurrency } from "@/lib/formatCurrency";

export default function Settings() {
  const { user, setRiskProfile, loadProfile, profile } = useAuthStore();
  const [autoDebit, setAutoDebit] = useState(true);
  const [notif, setNotif] = useState(true);
  const [darkOnly, setDarkOnly] = useState(true);

  const { data: referral } = useQuery({
    queryKey: ["referrals", "me"],
    queryFn: getReferralMe,
  });

  useEffect(() => {
    void loadProfile();
  }, [loadProfile]);

  return (
    <AppShell>
      <div className="mx-auto max-w-3xl p-6 md:p-10">
        <PageHeader title="Settings" subtitle="Account, money, and notifications." />

        <section className="rounded-3xl border border-border bg-card p-6">
          <h3 className="text-lg">Profile</h3>
          <div className="mt-4 grid gap-4 md:grid-cols-2">
            <Input label="Full name" defaultValue={user?.name} readOnly />
            <Input label="Email" defaultValue={user?.email} readOnly />
          </div>
          <p className="mt-3 text-xs text-muted-foreground">
            Effective penalty rate: {(profile?.effective_penalty_rate ?? user?.effectivePenaltyRate ?? 0.1) * 100}% ·
            Streak: {profile?.streak_count ?? user?.streakCount ?? 0} days
          </p>
        </section>

        <section className="mt-6 rounded-3xl border border-border bg-card p-6">
          <h3 className="text-lg">Referrals</h3>
          <p className="mt-2 text-sm text-muted-foreground">
            Share your code — friends get a discount when they sign up.
          </p>
          <div className="mt-4 flex items-center gap-3">
            <code className="rounded-xl bg-card-elevated px-4 py-2 font-mono text-lg text-coral">
              {referral?.code ?? "—"}
            </code>
            <span className="text-sm text-muted-foreground">
              {referral?.redemption_count ?? 0} redeemed · {referral?.discount_pct ?? 0}% off
            </span>
          </div>
        </section>

        <section className="mt-6 rounded-3xl border border-border bg-card p-6">
          <h3 className="text-lg">Risk profile</h3>
          <div className="mt-4 flex gap-2">
            {(["Conservative", "Moderate", "Aggressive"] as const).map((r) => (
              <button
                key={r}
                onClick={() => setRiskProfile(r)}
                className={`rounded-full border px-4 py-2 text-sm ${user?.riskProfile === r ? "border-coral bg-coral/15 text-coral" : "border-border text-muted-foreground hover:text-foreground"}`}
              >
                {r}
              </button>
            ))}
          </div>
          <p className="mt-3 text-xs text-muted-foreground">
            Invest threshold: {formatCurrency(profile?.invest_threshold ?? 500)}
          </p>
        </section>

        <section className="mt-6 rounded-3xl border border-border bg-card p-6 space-y-4">
          <h3 className="text-lg">Preferences</h3>
          <div className="flex items-center justify-between">
            <div><div className="text-sm">Auto-debit penalties</div><div className="text-xs text-muted-foreground">Move fines into portfolio automatically</div></div>
            <Toggle checked={autoDebit} onChange={setAutoDebit} />
          </div>
          <div className="flex items-center justify-between">
            <div><div className="text-sm">Push notifications</div><div className="text-xs text-muted-foreground">Penalty alerts &amp; weekly digest</div></div>
            <Toggle checked={notif} onChange={setNotif} />
          </div>
          <div className="flex items-center justify-between">
            <div><div className="text-sm">Dark theme only</div><div className="text-xs text-muted-foreground">Pebble looks best in the dark.</div></div>
            <Toggle checked={darkOnly} onChange={setDarkOnly} />
          </div>
        </section>

        <div className="mt-6">
          <Button variant="outline" onClick={() => void loadProfile()}>
            Refresh profile
          </Button>
        </div>
      </div>
    </AppShell>
  );
}
