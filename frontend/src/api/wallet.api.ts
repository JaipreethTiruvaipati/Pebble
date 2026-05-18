import { apiRequest } from "./client";
import type { WalletBalance, WalletLedgerEntry } from "@/types/api.types";

export function getWalletBalance() {
  return apiRequest<WalletBalance>("/wallet/balance");
}

export function getWalletLedger() {
  return apiRequest<WalletLedgerEntry[]>("/wallet/ledger");
}

export function topupWallet(amount: number) {
  return apiRequest<{ message: string; amount: number }>("/wallet/topup", {
    method: "POST",
    body: JSON.stringify({ amount }),
  });
}
