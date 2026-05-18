import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import * as penaltyApi from "@/api/penalty.api";

export function usePenalties(status?: string) {
  return useQuery({
    queryKey: ["penalties", status],
    queryFn: () => penaltyApi.listPenalties(status),
  });
}

export function useContestPenalty() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: penaltyApi.contestPenalty,
    onSuccess: () => qc.invalidateQueries({ queryKey: ["penalties"] }),
  });
}

export function useConfirmPenalty() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: penaltyApi.confirmPenaltyEarly,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["penalties"] });
      qc.invalidateQueries({ queryKey: ["wallet"] });
    },
  });
}
