import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import * as txApi from "@/api/transaction.api";

export function useTransactions(limit = 20) {
  return useQuery({
    queryKey: ["transactions", limit],
    queryFn: () => txApi.listTransactions(limit),
  });
}

export function useTransaction(id: string | undefined) {
  return useQuery({
    queryKey: ["transaction", id],
    queryFn: () => txApi.getTransaction(id!),
    enabled: !!id,
  });
}

export function useConfirmTransaction() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => txApi.confirmTransaction(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["transactions"] });
      qc.invalidateQueries({ queryKey: ["penalties"] });
      qc.invalidateQueries({ queryKey: ["wallet"] });
    },
  });
}
