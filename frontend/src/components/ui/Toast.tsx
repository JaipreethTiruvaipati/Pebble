import { AnimatePresence, motion } from "framer-motion";
import { CheckCircle2, AlertCircle, Info, X } from "lucide-react";
import { create } from "zustand";

type ToastType = "success" | "error" | "info";
interface ToastItem { id: string; message: string; type: ToastType }

interface ToastStore {
  toasts: ToastItem[];
  push: (message: string, type?: ToastType) => void;
  dismiss: (id: string) => void;
}

export const useToastStore = create<ToastStore>((set) => ({
  toasts: [],
  push: (message, type = "info") => {
    const id = Math.random().toString(36).slice(2);
    set((s) => ({ toasts: [...s.toasts, { id, message, type }] }));
    setTimeout(() => set((s) => ({ toasts: s.toasts.filter((t) => t.id !== id) })), 3500);
  },
  dismiss: (id) => set((s) => ({ toasts: s.toasts.filter((t) => t.id !== id) })),
}));

export function toast(message: string, type: ToastType = "info") {
  useToastStore.getState().push(message, type);
}

const icons = {
  success: <CheckCircle2 size={18} className="text-teal" />,
  error: <AlertCircle size={18} className="text-coral" />,
  info: <Info size={18} className="text-amber" />,
};

export function ToastContainer() {
  const { toasts, dismiss } = useToastStore();
  return (
    <div className="pointer-events-none fixed bottom-6 right-6 z-[100] flex flex-col gap-2">
      <AnimatePresence>
        {toasts.map((t) => (
          <motion.div
            key={t.id}
            initial={{ opacity: 0, x: 40, scale: 0.95 }}
            animate={{ opacity: 1, x: 0, scale: 1 }}
            exit={{ opacity: 0, x: 40, scale: 0.95 }}
            transition={{ duration: 0.25 }}
            className="pointer-events-auto flex items-center gap-3 rounded-xl bg-card-elevated border border-border px-4 py-3 shadow-2xl min-w-[260px]"
          >
            {icons[t.type]}
            <span className="text-sm flex-1">{t.message}</span>
            <button onClick={() => dismiss(t.id)} className="text-muted-foreground hover:text-foreground"><X size={14} /></button>
          </motion.div>
        ))}
      </AnimatePresence>
    </div>
  );
}
