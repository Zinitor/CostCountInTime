import { useState, useEffect, useCallback } from "react";
import type { KarmaState, KarmaUpdateResponse } from "../lib/karma";

const KARMA_STATE_KEY = "karma_state";

export function useKarma() {
  const [karma, setKarma] = useState<KarmaState | null>(null);
  const [loading, setLoading] = useState(true);

  // Загрузка состояния при монтировании
  useEffect(() => {
    fetchKarmaState();
  }, []);

  const fetchKarmaState = async () => {
    try {
      const res = await fetch("/api/karma/state");
      if (res.ok) {
        const data = await res.json();
        setKarma(data);
        localStorage.setItem(KARMA_STATE_KEY, JSON.stringify(data));
      } else {
        // Пробуем загрузить из localStorage
        const cached = localStorage.getItem(KARMA_STATE_KEY);
        if (cached) {
          setKarma(JSON.parse(cached));
        }
      }
    } catch {
      const cached = localStorage.getItem(KARMA_STATE_KEY);
      if (cached) {
        setKarma(JSON.parse(cached));
      }
    } finally {
      setLoading(false);
    }
  };

  const updateKarma = useCallback(
    async (
      price: number,
      hoursNeeded: number,
      itemName: string,
    ): Promise<KarmaUpdateResponse | null> => {
      try {
        const res = await fetch("/api/karma/update", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            price,
            hours_needed: hoursNeeded,
            item_name: itemName,
          }),
        });

        if (res.ok) {
          const data: KarmaUpdateResponse = await res.json();
          // Обновляем локальное состояние
          setKarma((prev) =>
            prev
              ? {
                  ...prev,
                  score: data.score,
                  level: data.level,
                  stats: data.stats,
                }
              : null,
          );
          return data;
        }
      } catch (err) {
        console.error("Failed to update karma:", err);
      }
      return null;
    },
    [],
  );

  const resetKarma = useCallback(async () => {
    try {
      const res = await fetch("/api/karma/reset", {
        method: "POST",
      });
      if (res.ok) {
        localStorage.removeItem(KARMA_STATE_KEY);
        setKarma(null);
        await fetchKarmaState();
        return true;
      }
    } catch (err) {
      console.error("Failed to reset karma:", err);
    }
    return false;
  }, []);

  const deleteKarmaEvent = useCallback(
    async (
      itemName: string,
      price: number,
      timestamp: number,
    ): Promise<boolean> => {
      try {
        const res = await fetch("/api/karma/delete_event", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            itemName,
            price,
            timestamp: Math.floor(timestamp / 1000), // конвертируем в секунды
          }),
        });
        if (res.ok) {
          // Обновляем состояние после удаления
          await fetchKarmaState();
          return true;
        }
      } catch (err) {
        console.error("Failed to delete karma event:", err);
      }
      return false;
    },
    [],
  );

  return {
    karma,
    loading,
    updateKarma,
    resetKarma,
    deleteKarmaEvent,
    refresh: fetchKarmaState,
  };
}
