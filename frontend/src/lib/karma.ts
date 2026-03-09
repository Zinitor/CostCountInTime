export type KarmaLevel =
  | "financial_god"
  | "wise_spender"
  | "normal_person"
  | "risky_buyer"
  | "financial_degen";

export interface KarmaStats {
  total_spent: number;
  total_hours: number;
  avg_purchase: number;
  impulse_count: number;
  best_streak: number;
  current_streak: number;
  purchases_today: number;
  spent_today: number;
}

export interface KarmaState {
  score: number;
  level: KarmaLevel;
  info: {
    level: KarmaLevel;
    min_score: number;
    max_score: number;
    display_name: string;
    icon: string;
    color: string;
  };
  stats: KarmaStats;
}

export interface KarmaUpdateResponse {
  score: number;
  level: KarmaLevel;
  score_delta: number;
  reasons: string[];
  stats: KarmaStats;
  comment: string; // Комментарий уровня кармы
  purchase_comment: string; // Комментарий покупки
}

export const karmaLevelInfo: Record<
  KarmaLevel,
  {
    displayName: string;
    icon: string;
    color: string;
    description: string;
  }
> = {
  financial_god: {
    displayName: "Financial God",
    icon: "😇",
    color: "#22c55e",
    description: "Ты мастер финансов! Так держать!",
  },
  wise_spender: {
    displayName: "Wise Spender",
    icon: "🧘",
    color: "#84cc16",
    description: "Разумный подход к тратам",
  },
  normal_person: {
    displayName: "Normal Person",
    icon: "😐",
    color: "#eab308",
    description: "Обычный человек, ничего особенного",
  },
  risky_buyer: {
    displayName: "Risky Buyer",
    icon: "😬",
    color: "#f97316",
    description: "Осторожно, опасные траты!",
  },
  financial_degen: {
    displayName: "Financial Degen",
    icon: "🤡",
    color: "#ef4444",
    description: "ТВОЙ КОШЕЛЁК В ОПАСНОСТИ!",
  },
};
