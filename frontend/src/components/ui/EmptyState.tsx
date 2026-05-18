import type { ReactNode } from "react";
import { Button } from "./Button";

export function EmptyState({ icon, title, description, action }: { icon?: ReactNode; title: string; description?: string; action?: { label: string; onClick: () => void } }) {
  return (
    <div className="flex flex-col items-center justify-center gap-3 rounded-xl border border-dashed border-border bg-card/40 p-10 text-center">
      {icon && <div className="text-muted-foreground">{icon}</div>}
      <h3 className="text-lg">{title}</h3>
      {description && <p className="max-w-sm text-sm text-muted-foreground">{description}</p>}
      {action && <Button size="sm" onClick={action.onClick}>{action.label}</Button>}
    </div>
  );
}
export default EmptyState;
