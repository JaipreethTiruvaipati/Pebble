import { Link, useLocation, useNavigate } from "@tanstack/react-router";
import { LayoutDashboard, ReceiptText, PieChart, BarChart3, Settings, History, Wallet, LogOut } from "lucide-react";
import type { ReactNode } from "react";
import { ROUTES } from "@/routes";
import { cn } from "@/lib/cn";
import { useAuthStore } from "@/stores/authStore";
import { ToastContainer } from "@/components/ui/Toast";

const nav = [
  { to: ROUTES.dashboard, label: "Dashboard", icon: LayoutDashboard },
  { to: ROUTES.logTransaction, label: "Log bill", icon: ReceiptText },
  { to: ROUTES.portfolio, label: "Portfolio", icon: PieChart },
  { to: ROUTES.insights, label: "Insights", icon: BarChart3 },
  { to: ROUTES.history, label: "History", icon: History },
  { to: ROUTES.settings, label: "Settings", icon: Settings },
];

export function AppShell({ children }: { children: ReactNode }) {
  const { user, logout } = useAuthStore();
  const { pathname } = useLocation();
  const navigate = useNavigate();

  return (
    <div className="flex min-h-screen bg-background text-foreground">
      <aside className="hidden md:flex w-60 shrink-0 flex-col border-r border-border bg-card/40 px-4 py-6">
        <Link to={ROUTES.landing} className="mb-8 flex items-center gap-2 px-2">
          <div className="grid h-9 w-9 place-items-center rounded-xl bg-coral text-primary-foreground">
            <Wallet size={18} />
          </div>
          <span className="font-display text-xl tracking-tight">Pebble</span>
        </Link>
        <nav className="flex flex-1 flex-col gap-1">
          {nav.map((n) => {
            const active = pathname === n.to || (n.to !== ROUTES.dashboard && pathname.startsWith(n.to));
            return (
              <Link
                key={n.to}
                to={n.to}
                className={cn(
                  "flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm transition-colors",
                  active ? "bg-coral/15 text-coral" : "text-muted-foreground hover:bg-card-elevated hover:text-foreground",
                )}
              >
                <n.icon size={16} />
                {n.label}
              </Link>
            );
          })}
        </nav>
        {user && (
          <div className="mt-6 flex items-center gap-3 rounded-xl border border-border bg-card p-3">
            <div className="grid h-9 w-9 place-items-center rounded-full bg-coral/20 text-sm font-semibold text-coral">{user.avatar}</div>
            <div className="flex-1 min-w-0">
              <div className="truncate text-sm">{user.name}</div>
              <div className="truncate text-[11px] text-muted-foreground">{user.riskProfile} profile</div>
            </div>
            <button
              type="button"
              className="text-muted-foreground hover:text-coral"
              onClick={() => {
                logout();
                navigate({ to: ROUTES.login });
              }}
            >
              <LogOut size={14} />
            </button>
          </div>
        )}
      </aside>
      <main className="flex-1 min-w-0">
        {children}
      </main>
      <ToastContainer />
    </div>
  );
}
export default AppShell;


