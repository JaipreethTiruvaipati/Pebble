import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import * as walletApi from "@/api/wallet.api";

export function useWallet() {
  const qc = useQueryClient();
  const balance = useQuery({
    queryKey: ["wallet", "balance"],
    queryFn: walletApi.getWalletBalance,
  });
  const ledger = useQuery({
    queryKey: ["wallet", "ledger"],
    queryFn: walletApi.getWalletLedger,
  });
  const topup = useMutation({
    mutationFn: (amount: number) => walletApi.topupWallet(amount),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["wallet"] });
    },
  });
  return { balance, ledger, topup };
}
