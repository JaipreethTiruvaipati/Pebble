import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { motion, AnimatePresence } from "framer-motion";
import { AppShell } from "@/components/layout/AppShell";
import { PageHeader } from "@/components/layout/PageHeader";
import { BillUploadZone } from "./components/BillUploadZone";
import { AnalysingSpinner } from "./components/AnalysingSpinner";
import { LineItemCard } from "./components/LineItemCard";
import { MerchantInput } from "./components/MerchantInput";
import { ConfirmAnimation } from "./components/ConfirmAnimation";
import { Button } from "@/components/ui/Button";
import { useDraftStore, type DraftLineItem } from "@/stores/draftStore";
import { formatCurrency } from "@/lib/formatCurrency";
import { calcPenalty } from "@/lib/calcPenalty";
import { toast } from "@/components/ui/Toast";
import { ROUTES } from "@/routes";
import * as txApi from "@/api/transaction.api";
import { useAuthStore } from "@/stores/authStore";

type Step = "upload" | "preview" | "analysing" | "review" | "confirmed";

export default function LogTransaction() {
  const [step, setStep] = useState<Step>("upload");
  const [merchant, setMerchant] = useState("");
  const [file, setFile] = useState<File | null>(null);
  const [txId, setTxId] = useState<string | null>(null);
  const { draft, setDraft } = useDraftStore();
  const navigate = useNavigate();
  const profile = useAuthStore((s) => s.profile);
  const rate = profile?.effective_penalty_rate ?? 0.1;
  const threshold = profile?.penalty_threshold ?? 50;

  const handleFile = (f: File) => {
    setFile(f);
    setStep("preview");
  };

  const startAnalysis = async () => {
    if (!merchant.trim()) {
      toast("Enter a merchant name", "error");
      return;
    }
    setStep("analysing");
    try {
      const form = new FormData();
      form.append("merchant", merchant);
      form.append("total_amount", "1000");
      if (file) form.append("receipt", file);

      const res = await txApi.uploadBill(form);
      const id = res.transaction_id;
      setTxId(id);

      const detail = await txApi.pollTransactionUntilScored(id, 20, 1500);
      const items: DraftLineItem[] = detail.line_items.map((li) => ({
        id: li.id,
        name: li.name,
        amount: li.amount,
        category: li.category || "other",
        score: li.impulse_score,
        reasoning: li.reasoning || "Scored by Pebble AI",
      }));
      setDraft({ merchant: detail.merchant, date: detail.logged_at, items });
      setStep("review");
    } catch (e) {
      toast(e instanceof Error ? e.message : "Analysis failed", "error");
      setStep("preview");
    }
  };

  const items = draft?.items ?? [];
  const totalAmount = items.reduce((s, i) => s + i.amount, 0);
  const totalPenalty = items.reduce((s, i) => s + calcPenalty(i.amount, i.score, rate, threshold), 0);

  const confirm = async () => {
    if (!txId) return;
    try {
      await txApi.confirmTransaction(txId);
      setStep("confirmed");
      toast(`Penalty of ${formatCurrency(totalPenalty)} queued`, "success");
      setTimeout(() => {
        setDraft(null);
        navigate({ to: ROUTES.dashboard });
      }, 1800);
    } catch {
      toast("Confirm failed", "error");
    }
  };

  return (
    <AppShell>
      <div className="mx-auto max-w-4xl p-6 md:p-10">
        <PageHeader title="Log a transaction" subtitle="We extract items and score each one for impulsiveness." />

        <AnimatePresence mode="wait">
          {step === "upload" && (
            <motion.div key="upload" initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
              <BillUploadZone onFile={handleFile} />
            </motion.div>
          )}

          {step === "preview" && (
            <motion.div key="preview" initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }} exit={{ opacity: 0 }}>
              <div className="rounded-3xl border border-border bg-card p-8">
                <h3 className="text-xl">Confirm details</h3>
                <p className="mt-1 text-sm text-muted-foreground">We&apos;ll scan this bill via the scoring pipeline.</p>
                <div className="mt-6 grid gap-5 md:grid-cols-2">
                  <MerchantInput value={merchant} onChange={setMerchant} />
                  <div className="rounded-xl border border-dashed border-border bg-card-elevated p-4">
                    <div className="text-xs text-muted-foreground">Receipt</div>
                    <div className="mt-2 text-sm">{file?.name ?? "No file"}</div>
                  </div>
                </div>
                <div className="mt-6 flex gap-3">
                  <Button variant="outline" onClick={() => setStep("upload")}>
                    Back
                  </Button>
                  <Button onClick={startAnalysis}>Analyse bill →</Button>
                </div>
              </div>
            </motion.div>
          )}

          {step === "analysing" && (
            <motion.div key="anal" initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
              <AnalysingSpinner onDone={() => {}} />
            </motion.div>
          )}

          {step === "review" && (
            <motion.div
              key="review"
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0 }}
              className="space-y-4"
            >
              <div className="flex items-center justify-between rounded-2xl border border-border bg-card px-5 py-4">
                <div>
                  <div className="text-xs text-muted-foreground">{merchant}</div>
                  <div className="mt-1 font-mono text-lg">
                    {formatCurrency(totalAmount)} <span className="text-xs text-muted-foreground">total</span>
                  </div>
                </div>
                <div className="text-right">
                  <div className="text-xs text-muted-foreground">{items.length} line items</div>
                  <div className="font-mono text-sm text-coral">−{formatCurrency(totalPenalty)} penalty</div>
                </div>
              </div>

              <div className="space-y-3">
                {items.map((it, i) => (
                  <LineItemCard key={it.id} item={it} index={i} />
                ))}
              </div>

              <div className="sticky bottom-4 mt-6 rounded-2xl border border-coral/40 bg-card/95 px-5 py-4 backdrop-blur">
                <div className="flex items-center justify-between gap-4">
                  <div>
                    <div className="text-xs uppercase tracking-widest text-muted-foreground">Total penalty</div>
                    <div className="mt-0.5 font-mono text-3xl text-coral">{formatCurrency(totalPenalty)}</div>
                  </div>
                  <div className="flex gap-2">
                    <Button variant="outline" onClick={() => setStep("upload")}>
                      Review
                    </Button>
                    <Button variant="coral" onClick={confirm}>
                      Accept &amp; invest
                    </Button>
                  </div>
                </div>
              </div>
            </motion.div>
          )}

          {step === "confirmed" && (
            <motion.div key="done" initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}>
              <ConfirmAnimation />
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </AppShell>
  );
}
