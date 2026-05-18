import { apiRequest } from "./client";
import type { TransactionDetail, TransactionSummary } from "@/types/api.types";

export function createTransaction(merchant: string, totalAmount: number) {
  return apiRequest<{ transaction_id: string; status: string }>("/transactions", {
    method: "POST",
    body: JSON.stringify({ merchant, total_amount: totalAmount }),
  });
}

export function uploadBill(form: FormData) {
  return apiRequest<{ transaction_id: string; message?: string }>("/transactions/bill", {
    method: "POST",
    body: form,
  });
}

export function listTransactions(limit = 20) {
  return apiRequest<TransactionSummary[]>(`/transactions?limit=${limit}`);
}

export function getTransaction(id: string) {
  return apiRequest<TransactionDetail>(`/transactions/${id}`);
}

export function confirmTransaction(id: string) {
  return apiRequest<{
    transaction_id: string;
    penalties_created: number;
    total_penalty_queued: number;
    status: string;
  }>(`/transactions/${id}/confirm`, { method: "POST" });
}

export function overrideLineItemScore(lineItemId: string, overrideScore: number) {
  return apiRequest(`/line-items/${lineItemId}/score`, {
    method: "PUT",
    body: JSON.stringify({ override_score: overrideScore }),
  });
}

export async function pollTransactionUntilScored(
  id: string,
  maxAttempts = 30,
  intervalMs = 2000,
): Promise<TransactionDetail> {
  for (let i = 0; i < maxAttempts; i++) {
    const tx = await getTransaction(id);
    if (tx.status === "scored" && tx.line_items.length > 0) return tx;
    await new Promise((r) => setTimeout(r, intervalMs));
  }
  return getTransaction(id);
}
