import { motion } from "framer-motion";
import { Check } from "lucide-react";

export function ConfirmAnimation() {
  return (
    <div className="mx-auto flex max-w-md flex-col items-center rounded-3xl border border-teal/30 bg-teal/10 p-12 text-center">
      <motion.div
        initial={{ scale: 0 }} animate={{ scale: 1 }} transition={{ type: "spring", stiffness: 200, damping: 16 }}
        className="grid h-20 w-20 place-items-center rounded-full bg-teal text-background"
      >
        <Check size={36} strokeWidth={3} />
      </motion.div>
      <motion.h3
        initial={{ opacity: 0, y: 8 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.3 }}
        className="mt-6 text-2xl"
      >
        Bill logged
      </motion.h3>
      <motion.p
        initial={{ opacity: 0 }} animate={{ opacity: 1 }} transition={{ delay: 0.5 }}
        className="mt-2 text-sm text-muted-foreground"
      >
        Penalties routed to your portfolio.
      </motion.p>
    </div>
  );
}
export default ConfirmAnimation;
