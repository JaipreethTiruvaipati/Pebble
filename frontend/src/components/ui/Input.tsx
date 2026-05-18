import { cn } from "@/lib/cn";
import { forwardRef, type InputHTMLAttributes } from "react";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  hint?: string;
  error?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(({ className, label, hint, error, ...props }, ref) => {
  return (
    <label className="flex flex-col gap-1.5">
      {label && <span className="text-xs font-medium text-muted-foreground">{label}</span>}
      <input
        ref={ref}
        className={cn(
          "h-11 w-full rounded-lg bg-input border border-border px-3.5 text-sm text-foreground placeholder:text-muted-foreground/60 outline-none transition-colors focus:border-coral focus:ring-2 focus:ring-coral/30",
          error && "border-coral focus:border-coral focus:ring-coral/40",
          className,
        )}
        {...props}
      />
      {(hint || error) && (
        <span className={cn("text-[11px]", error ? "text-coral" : "text-muted-foreground")}>{error ?? hint}</span>
      )}
    </label>
  );
});
Input.displayName = "Input";
export default Input;

