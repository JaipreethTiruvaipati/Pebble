import { useState } from "react";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { Toggle } from "@/components/ui/Toggle";
import { Input } from "@/components/ui/Input";
import { useAuthStore } from "@/stores/authStore";

export default function Settings() {
  const { user, setRiskProfile } = useAuthStore();
  const [autoDebit, setAutoDebit] = useState(true);
  const [notif, setNotif] = useState(true);
  const [darkOnly, setDarkOnly] = useState(true);

  return (
    <AppShell>
      <div className="mx-auto max-w-3xl p-6 md:p-10">
        <PageHeader title="Settings" subtitle="Account, money, and notifications." />

        <section className="rounded-3xl border border-border bg-card p-6">
          <h3 className="text-lg">Profile</h3>
          <div className="mt-4 grid gap-4 md:grid-cols-2">
            <Input label="Full name" defaultValue={user?.name} />
            <Input label="Email" defaultValue={user?.email} />
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
      </div>
    </AppShell>
  );
}
