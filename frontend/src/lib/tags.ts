export interface Tag {
  id: string;
  name: string;
  colorStart: string;
  colorEnd: string;
}

export const DEFAULT_TAGS: Tag[] = [
  {
    id: "food",
    name: "Еда",
    colorStart: "#ef4444",
    colorEnd: "#dc2626",
  },
  {
    id: "entertainment",
    name: "Развлечения",
    colorStart: "#f59e0b",
    colorEnd: "#d97706",
  },
  {
    id: "transport",
    name: "Транспорт",
    colorStart: "#3b82f6",
    colorEnd: "#2563eb",
  },
  {
    id: "shopping",
    name: "Покупки",
    colorStart: "#ec4899",
    colorEnd: "#db2777",
  },
  {
    id: "tech",
    name: "Техника",
    colorStart: "#8b5cf6",
    colorEnd: "#7c3aed",
  },
];
