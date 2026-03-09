import { useState, useEffect } from "react";
import { DollarSign, ShoppingCart, TrendingUp, Award } from "lucide-react";
import { useKarma } from "./hooks/useKarma";
import { KarmaBadge } from "./components/KarmaBadge";
import { KarmaModal } from "./components/KarmaModal";
import { KarmaBackground } from "./components/KarmaBackground";
import { KarmaAlert } from "./components/KarmaAlert";
import { PurchaseCard } from "./components/PurchaseCard";
import { KarmaResetAnimation } from "./components/KarmaResetAnimation";

interface Purchase {
  id: number;
  itemName: string;
  price: number;
  hours: number;
  days: number;
  formatted: string;
  comment: string;
  timestamp: number;
  purchaseComment?: string;
}

type SalaryPeriod = "hour" | "day" | "month";

const periodLabels: Record<SalaryPeriod, string> = {
  hour: "час",
  day: "день",
  month: "месяц",
};

const STORAGE_KEY = "costcount_purchases";
const SALARY_KEY = "costcount_salary";
const SALARY_PERIOD_KEY = "costcount_salary_period";
const EXPIRY_HOURS = 24; // 1 день

function App() {
  const [salary, setSalary] = useState<number>(50000);
  const [salaryPeriod, setSalaryPeriod] = useState<SalaryPeriod>("month");
  const [itemName, setItemName] = useState<string>("");
  const [price, setPrice] = useState<number>(0);
  const [purchases, setPurchases] = useState<Purchase[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<string>("");

  // Karma
  const {
    karma,
    loading: karmaLoading,
    updateKarma,
    resetKarma,
    deleteKarmaEvent,
  } = useKarma();
  const [showKarmaModal, setShowKarmaModal] = useState(false);
  const [lastKarmaComment, setLastKarmaComment] = useState<string>("");
  const [karmaAlert, setKarmaAlert] = useState<{
    show: boolean;
    level: any;
    score: number;
    scoreDelta: number;
    purchaseComment: string;
  } | null>(null);
  const [justAddedIds, setJustAddedIds] = useState<Set<number>>(new Set());
  const [showResetAnimation, setShowResetAnimation] = useState(false);

  // Загрузка данных при монтировании
  useEffect(() => {
    const now = Date.now();
    const expiryMs = EXPIRY_HOURS * 60 * 60 * 1000;

    // Загружаем зарплату
    const savedSalary = localStorage.getItem(SALARY_KEY);
    const savedPeriod = localStorage.getItem(SALARY_PERIOD_KEY);
    if (savedSalary) setSalary(Number(savedSalary));
    if (savedPeriod) setSalaryPeriod(savedPeriod as SalaryPeriod);

    // Загружаем покупки, фильтруем просроченные
    const savedPurchases = localStorage.getItem(STORAGE_KEY);
    if (savedPurchases) {
      try {
        const parsed: Purchase[] = JSON.parse(savedPurchases);
        const valid = parsed.filter((p) => now - p.timestamp < expiryMs);
        setPurchases(valid);
        // Сохраняем обратно только валидные
        if (valid.length !== parsed.length) {
          localStorage.setItem(STORAGE_KEY, JSON.stringify(valid));
        }
      } catch {
        localStorage.removeItem(STORAGE_KEY);
      }
    }
  }, []);

  // Сохранение покупок при изменении
  useEffect(() => {
    if (purchases.length > 0) {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(purchases));
    } else {
      localStorage.removeItem(STORAGE_KEY);
    }
  }, [purchases]);

  // Сохранение зарплаты при изменении
  useEffect(() => {
    localStorage.setItem(SALARY_KEY, String(salary));
  }, [salary]);

  useEffect(() => {
    localStorage.setItem(SALARY_PERIOD_KEY, salaryPeriod);
  }, [salaryPeriod]);

  // Пересчёт всех покупок при изменении зарплаты
  useEffect(() => {
    if (purchases.length > 0) {
      recalculateAllPurchases();
    }
  }, [salary, salaryPeriod]);

  const recalculateAllPurchases = async () => {
    const updated = await Promise.all(
      purchases.map(async (p) => {
        try {
          const res = await fetch("/api/calculate", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
              salary,
              salaryPeriod,
              price: p.price,
              itemName: p.itemName,
            }),
          });
          if (res.ok) {
            const data = await res.json();
            return {
              ...p,
              hours: data.hours,
              days: data.days,
              formatted: data.formatted,
              comment: data.comment,
            };
          }
        } catch {
          // Игнорируем ошибки пересчёта
        }
        return p;
      }),
    );
    setPurchases(updated);
  };

  const calculate = async () => {
    if (price <= 0) {
      setError("Введите корректную цену");
      return;
    }

    setLoading(true);
    setError("");

    try {
      const res = await fetch("/api/calculate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          salary,
          salaryPeriod,
          price,
          itemName: itemName || "Без названия",
        }),
      });

      if (!res.ok) throw new Error("Ошибка расчёта");

      const data = await res.json();

      // Обновляем карму
      const karmaResult = await updateKarma(
        price,
        data.hours,
        itemName || "Без названия",
      );
      if (karmaResult) {
        // Показывем комментарий покупки (не уровня кармы!)
        setLastKarmaComment(
          karmaResult.purchase_comment || karmaResult.comment,
        );
        // Показываем KarmaAlert
        setKarmaAlert({
          show: true,
          level: karmaResult.level,
          score: karmaResult.score,
          scoreDelta: karmaResult.score_delta,
          purchaseComment: karmaResult.purchase_comment || karmaResult.comment,
        });
      }

      // Добавляем новую покупку в НАЧАЛО списка
      const newId = Date.now();
      setPurchases((prev) => [
        {
          id: newId,
          itemName: data.itemName,
          price: data.price,
          hours: data.hours,
          days: data.days,
          formatted: data.formatted,
          comment: data.comment,
          timestamp: Date.now(),
          purchaseComment: karmaResult?.purchase_comment || "",
        },
        ...prev,
      ]);

      // Добавляем ID в список только что добавленных для анимации
      setJustAddedIds((prev) => new Set(prev).add(newId));
      // Удаляем из списка анимации через 500мс
      setTimeout(() => {
        setJustAddedIds((prev) => {
          const next = new Set(prev);
          next.delete(newId);
          return next;
        });
      }, 500);

      // Очищаем только цену, название оставляем
      setPrice(0);
    } catch {
      setError("Не удалось выполнить расчёт");
    } finally {
      setLoading(false);
    }
  };

  const removePurchase = async (id: number) => {
    const purchase = purchases.find((p) => p.id === id);
    if (purchase) {
      // Удаляем событие кармы, связанное с этой покупкой
      await deleteKarmaEvent(
        purchase.itemName,
        purchase.price,
        purchase.timestamp,
      );
    }
    setPurchases((prev) => prev.filter((p) => p.id !== id));
  };

  const handleResetKarma = () => {
    // Запускаем анимацию сброса кармы
    setShowResetAnimation(true);
  };

  const handleResetAnimationComplete = async () => {
    // После завершения анимации реально сбрасываем карму
    const success = await resetKarma();
    if (success) {
      setLastKarmaComment("Карма сброшена! Начинаем с чистого листа! 😇");
    }
    setShowResetAnimation(false);
  };

  return (
    <>
      {/* Анимированный фон в зависимости от кармы */}
      <KarmaBackground karmaLevel={karma?.level || null} />

      <div className="container mx-auto px-4 py-8 max-w-2xl relative z-10">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="flex justify-between items-center mb-4">
            <h1 className="text-4xl font-bold text-white">
              💰 Cost Count In Time
            </h1>
            {!karmaLoading && karma && (
              <KarmaBadge
                karma={karma}
                onClick={() => setShowKarmaModal(true)}
              />
            )}
          </div>
          <p className="text-purple-200">
            Узнайте, сколько нужно работать для покупки мечты
          </p>
        </div>

        {/* Последний комментарий кармы */}
        {lastKarmaComment && (
          <div className="mb-6 bg-gradient-to-r from-purple-500/20 to-pink-500/20 backdrop-blur-lg rounded-xl p-4 border border-purple-500/30">
            <div className="flex items-start gap-3">
              <Award className="w-6 h-6 text-purple-400 flex-shrink-0 mt-0.5" />
              <div>
                <div className="text-purple-200 text-sm font-medium mb-1">
                  Финансовая карма говорит:
                </div>
                <div className="text-white italic">{lastKarmaComment}</div>
              </div>
              <button
                onClick={() => setLastKarmaComment("")}
                className="text-purple-400 hover:text-white ml-auto"
              >
                ✕
              </button>
            </div>
          </div>
        )}

        {/* Salary Card */}
        <div className="bg-white/10 backdrop-blur-lg rounded-2xl p-6 mb-6 border border-white/20">
          <div className="flex items-center gap-3 mb-4">
            <DollarSign className="w-6 h-6 text-green-400" />
            <h2 className="text-xl font-semibold text-white">Ваша зарплата</h2>
          </div>

          <div className="space-y-4">
            <div>
              <label className="block text-purple-200 text-sm mb-2">
                Сумма (₽)
              </label>
              <input
                type="number"
                value={salary}
                onChange={(e) => setSalary(Number(e.target.value))}
                className="w-full bg-white/5 border border-white/20 rounded-xl px-4 py-3 text-white text-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                placeholder="50000"
              />
            </div>

            <div>
              <label className="block text-purple-200 text-sm mb-2">
                Период
              </label>
              <div className="flex gap-2">
                {(["hour", "day", "month"] as SalaryPeriod[]).map((period) => (
                  <button
                    key={period}
                    onClick={() => setSalaryPeriod(period)}
                    className={`flex-1 py-2 px-4 rounded-xl font-medium transition-all ${
                      salaryPeriod === period
                        ? "bg-purple-500 text-white"
                        : "bg-white/5 text-purple-200 hover:bg-white/10"
                    }`}
                  >
                    в {periodLabels[period]}
                  </button>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Purchase Input Card */}
        <div className="bg-white/10 backdrop-blur-lg rounded-2xl p-6 mb-6 border border-white/20">
          <div className="flex items-center gap-3 mb-4">
            <ShoppingCart className="w-6 h-6 text-pink-400" />
            <h2 className="text-xl font-semibold text-white">
              Что хотите купить?
            </h2>
          </div>

          <div className="space-y-4">
            <div>
              <label className="block text-purple-200 text-sm mb-2">
                Название{" "}
                <span className="text-purple-400">(необязательно)</span>
              </label>
              <input
                type="text"
                value={itemName}
                onChange={(e) => setItemName(e.target.value)}
                className="w-full bg-white/5 border border-white/20 rounded-xl px-4 py-3 text-white text-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                placeholder="Например: iPhone 15"
              />
            </div>

            <div>
              <label className="block text-purple-200 text-sm mb-2">
                Цена (₽)
              </label>
              <input
                type="number"
                value={price || ""}
                onChange={(e) => setPrice(Number(e.target.value))}
                className="w-full bg-white/5 border border-white/20 rounded-xl px-4 py-3 text-white text-lg focus:outline-none focus:ring-2 focus:ring-purple-500"
                placeholder="120000"
              />
            </div>

            {error && (
              <div className="text-red-400 text-sm bg-red-500/10 px-4 py-2 rounded-lg">
                {error}
              </div>
            )}

            <button
              onClick={calculate}
              disabled={loading}
              className="w-full bg-gradient-to-r from-purple-500 to-pink-500 hover:from-purple-600 hover:to-pink-600 disabled:opacity-50 text-white font-semibold py-4 rounded-xl transition-all transform hover:scale-[1.02] active:scale-[0.98]"
            >
              {loading ? "Считаем..." : "Рассчитать"}
            </button>
          </div>
        </div>

        {/* Results */}
        {purchases.length > 0 && (
          <div className="space-y-4">
            <div className="flex items-center gap-3 mb-4">
              <TrendingUp className="w-6 h-6 text-yellow-400" />
              <h2 className="text-xl font-semibold text-white">Результаты</h2>
            </div>

            {purchases.map((purchase) => (
              <PurchaseCard
                key={purchase.id}
                id={purchase.id}
                itemName={purchase.itemName}
                price={purchase.price}
                hours={purchase.hours}
                formatted={purchase.formatted}
                comment={purchase.comment}
                purchaseComment={purchase.purchaseComment}
                onRemove={removePurchase}
                isNew={justAddedIds.has(purchase.id)}
              />
            ))}
          </div>
        )}

        {/* Karma Modal */}
        {showKarmaModal && karma && (
          <KarmaModal
            karma={karma}
            onClose={() => setShowKarmaModal(false)}
            onReset={handleResetKarma}
          />
        )}

        {/* Karma Alert */}
        {karmaAlert?.show && karma && (
          <KarmaAlert
            level={karmaAlert.level}
            score={karmaAlert.score}
            scoreDelta={karmaAlert.scoreDelta}
            purchaseComment={karmaAlert.purchaseComment}
            onClose={() => setKarmaAlert(null)}
          />
        )}

        {/* Karma Reset Animation */}
        {showResetAnimation && (
          <KarmaResetAnimation onComplete={handleResetAnimationComplete} />
        )}
      </div>
    </>
  );
}

export default App;
