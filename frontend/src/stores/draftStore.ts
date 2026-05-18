import { create } from "zustand";

export interface DraftLineItem {
  id: string;
  name: string;
  amount: number;
  category: string;
  score: number;
  reasoning: string;
}

export interface BillDraft {
  merchant: string;
  date: string;
  items: DraftLineItem[];
  imageUrl?: string;
}

interface DraftStore {
  draft: BillDraft | null;
  setDraft: (d: BillDraft | null) => void;
  updateItem: (id: string, patch: Partial<DraftLineItem>) => void;
  removeItem: (id: string) => void;
}

export const useDraftStore = create<DraftStore>((set) => ({
  draft: null,
  setDraft: (draft) => set({ draft }),
  updateItem: (id, patch) =>
    set((s) => (s.draft ? { draft: { ...s.draft, items: s.draft.items.map((i) => (i.id === id ? { ...i, ...patch } : i)) } } : s)),
  removeItem: (id) => set((s) => (s.draft ? { draft: { ...s.draft, items: s.draft.items.filter((i) => i.id !== id) } } : s)),
}));
