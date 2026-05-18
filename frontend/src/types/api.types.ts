export interface AuthResponse {
  token: string;
  user_id: string;
}

export interface UserProfile {
  id: string;
  email: string;
  phone: string;
  risk_profile: string;
  penalty_rate: number;
  effective_penalty_rate: number;
  penalty_threshold: number;
  invest_threshold: number;
  consent_hours: number;
  streak_count: number;
  streak_last_updated?: string | null;
  streak_discount_pct: number;
}

export interface WalletBalance {
  balance: number;
  pending_total: number;
  invested_total: number;
}

export interface WalletLedgerEntry {
  id: string;
  type: string;
  amount: number;
  balance_after: number;
  created_at: string;
}

export interface LineItemDetail {
  id: string;
  name: string;
  amount: number;
  impulse_score: number;
  category: string;
  reasoning: string;
  user_overridden: boolean;
}

export interface TransactionDetail {
  id: string;
  merchant: string;
  total_amount: number;
  status: string;
  logged_at: string;
  line_items: LineItemDetail[];
}

export interface TransactionSummary {
  id: string;
  merchant: string;
  total_amount: number;
  status: string;
  logged_at: string;
  avg_score: number;
  total_penalty: number;
}

export interface PenaltyRow {
  id: string;
  line_item_id: string;
  amount: number;
  status: string;
  expires_at: string;
  merchant?: string;
  item_name?: string;
}

export interface PendingPenaltyBanner {
  penalty_id: string;
  amount: number;
  source: string;
  expires_at: string;
}

export interface PortfolioResponse {
  total_invested: number;
  equity_value: number;
  gold_value: number;
  bond_value: number;
  gain_pct: number;
  allocation_pct: Record<string, number>;
}

export interface Investment {
  id: string;
  user_id: string;
  asset_class: string;
  amount: number;
  units: number;
  nav_at_purchase: number;
  status: string;
  trigger_type?: string;
  broker_ref?: string;
  created_at: string;
}

export interface InvestmentListResponse {
  investments: Investment[];
  total: number;
}

export interface MarketSignal {
  asset_class: string;
  indicator: string;
  value: number;
  action: string;
  timestamp: string;
}

export interface MarketSignalResponse {
  signals: MarketSignal[];
  composite_score: number;
  updated_at: string;
}

export interface CategorySpend {
  category: string;
  amount: number;
  pct: number;
}

export interface WeeklyDigest {
  week_start: string;
  week_end: string;
  total_spend: number;
  impulse_pct: number;
  avg_impulse_score: number;
  top_categories: CategorySpend[];
  trend_vs_last_week_pct: number;
}

export interface BenchmarkResult {
  user_impulse_pct: number;
  cohort_impulse_pct: number;
  saved_vs_cohort_pct: number;
  cohort_label: string;
  sample_size: number;
}

export interface ReferralStats {
  code: string;
  redemption_count: number;
  discount_pct: number;
}

export interface ApiError {
  error: string;
  code: string;
}
