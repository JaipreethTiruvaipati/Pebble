export type TriggerType = 'threshold' | 'opportunity' | 'time' | string;

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

export interface PortfolioResponse {
  total_invested: number;
  equity_value: number;
  gold_value: number;
  bond_value: number;
  gain_pct: number;
  allocation_pct: Record<string, number>;
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

export interface InvestmentListResponse {
  investments: Investment[];
  total: number;
}
