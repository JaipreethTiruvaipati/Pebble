// Central mock dataset for the Pebble UI inside firebase.ts to match structure.
import type { DraftLineItem } from "@/stores/draftStore";

export interface Transaction {
  id: string;
  merchant: string;
  category: string;
  amount: number;
  score: number;
  penalty: number;
  date: string;
}

export const wallet = {
  balance: 124380,
  invested: 86400,
  investedThisMonth: 12450,
  impulseSavesCount: 17,
  impulseSavedAmount: 8920,
  monthlyTarget: 25000,
};

export const pendingPenalty = {
  amount: 1240,
  source: "Sneakers — Nike Air Max",
  expiresAt: new Date(Date.now() + 18 * 3600 * 1000 + 23 * 60 * 1000).toISOString(),
};

export const recentTransactions: Transaction[] = [
  { id: "t1", merchant: "Zomato", category: "Food", amount: 680, score: 72, penalty: 92, date: "2026-05-17T20:14:00" },
  { id: "t2", merchant: "Amazon", category: "Shopping", amount: 4299, score: 88, penalty: 612, date: "2026-05-16T10:02:00" },
  { id: "t3", merchant: "BMTC Metro", category: "Travel", amount: 60, score: 8, penalty: 0, date: "2026-05-16T08:31:00" },
  { id: "t4", merchant: "BigBasket", category: "Groceries", amount: 2140, score: 22, penalty: 18, date: "2026-05-15T18:55:00" },
  { id: "t5", merchant: "Starbucks", category: "Food", amount: 420, score: 64, penalty: 48, date: "2026-05-15T11:20:00" },
  { id: "t6", merchant: "Apple Store", category: "Tech", amount: 18900, score: 92, penalty: 2820, date: "2026-05-14T17:42:00" },
];

export const portfolioTimeline = [
  { day: "W1", value: 64000 }, { day: "W2", value: 66200 }, { day: "W3", value: 65800 },
  { day: "W4", value: 68900 }, { day: "W5", value: 71200 }, { day: "W6", value: 73400 },
  { day: "W7", value: 76800 }, { day: "W8", value: 79100 }, { day: "W9", value: 81000 },
  { day: "W10", value: 82400 }, { day: "W11", value: 84600 }, { day: "W12", value: 86400 },
];

export const holdings = [
  { name: "Equity (NIFTY 50 ETF)", value: 48000, allocation: 55.6, change: 12.4 },
  { name: "Sovereign Gold Bonds", value: 17280, allocation: 20.0, change: 8.1 },
  { name: "Corporate Bonds", value: 12960, allocation: 15.0, change: 4.2 },
  { name: "Liquid Fund", value: 8160, allocation: 9.4, change: 1.1 },
];

export const allocation = [
  { name: "Equity", value: 55.6, color: "#E94560" },
  { name: "Gold", value: 20, color: "#F5A623" },
  { name: "Bonds", value: 15, color: "#00D4AA" },
  { name: "Liquid", value: 9.4, color: "#8B8FA8" },
];

export const categoryBreakdown = [
  { category: "Food", spend: 8420, score: 58 },
  { category: "Shopping", spend: 14290, score: 82 },
  { category: "Travel", spend: 3120, score: 18 },
  { category: "Groceries", spend: 6480, score: 22 },
  { category: "Entertainment", spend: 4210, score: 71 },
  { category: "Tech", spend: 18900, score: 92 },
];

export const topMerchants = [
  { merchant: "Amazon", spend: 12640, score: 84 },
  { merchant: "Zomato", spend: 4820, score: 68 },
  { merchant: "Apple", spend: 18900, score: 92 },
  { merchant: "Starbucks", spend: 2140, score: 62 },
  { merchant: "BigBasket", spend: 6480, score: 24 },
];

export const weeklyDigest = {
  thisWeek: 12640,
  lastWeek: 16820,
  topCategory: "Shopping",
  saved: 1240,
};

export const peerBenchmark = {
  you: 64,
  peers: 78,
  cohort: "26-32 Moderate",
};

export const streak = { current: 8, best: 21 };

export const marketPulse = { score: 72, label: "Bullish", note: "NIFTY 50 entering accumulation zone — moderate buys favored." };

export const mockLineItems: DraftLineItem[] = [
  { id: "li1", name: "Nike Air Max 270", amount: 12990, category: "Shopping", score: 88, reasoning: "Late-night browse, 3rd pair this quarter. Discretionary apparel." },
  { id: "li2", name: "AirPods Pro Case", amount: 2199, category: "Tech", score: 74, reasoning: "Accessory purchase shortly after recent gadget buy." },
  { id: "li3", name: "Protein Bar pack", amount: 480, category: "Health", score: 22, reasoning: "Recurring health item, aligned with goals." },
  { id: "li4", name: "Phone charger", amount: 1290, category: "Tech", score: 38, reasoning: "Replacement item — useful but not urgent." },
];

// Firebase Config Placeholder (if needed in the future)
export const db = null;
export const auth = null;
