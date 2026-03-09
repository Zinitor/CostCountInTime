import { useState, useEffect } from "react";
import { Clock, DollarSign, ShoppingCart, TrendingUp } from "lucide-react";

interface Purchase {
  id: number;
  itemName: string;
  price: number;
  hours: number;
  days: number;
  formatted: string;
  comment: string;
  timestamp: number;
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

      // Добавляем новую покупку в НАЧАЛО списка
      setPurchases((prev) => [
        {
          id: Date.now(),
          itemName: data.itemName,
          price: data.price,
          hours: data.hours,
          days: data.days,
          formatted: data.formatted,
          comment: data.comment,
          timestamp: Date.now(),
        },
        ...prev,
      ]);

      // Очищаем только цену, название оставляем
      setPrice(0);
    } catch {
      setError("Не удалось выполнить расчёт");
    } finally {
      setLoading(false);
    }
  };

  const removePurchase = (id: number) => {
    setPurchases((prev) => prev.filter((p) => p.id !== id));
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      {/* Header */}
      <div className="text-center mb-8">
        <h1 className="text-4xl font-bold text-white mb-2">
          💰 Cost Count In Time
        </h1>
        <p className="text-purple-200">
          Узнайте, сколько нужно работать для покупки мечты
        </p>
      </div>

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
            <label className="block text-purple-200 text-sm mb-2">Период</label>
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
              Название <span className="text-purple-400">(необязательно)</span>
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
            <div
              key={purchase.id}
              className="bg-gradient-to-r from-purple-500/20 to-pink-500/20 backdrop-blur-lg rounded-2xl p-6 border border-white/20"
            >
              <div className="flex justify-between items-start mb-3">
                <div>
                  <h3 className="text-lg font-semibold text-white">
                    {purchase.itemName}
                  </h3>
                  <p className="text-purple-200">
                    {purchase.price.toLocaleString()} ₽
                  </p>
                </div>
                <button
                  onClick={() => removePurchase(purchase.id)}
                  className="text-purple-300 hover:text-white transition-colors"
                >
                  ✕
                </button>
              </div>

              <div className="flex items-center gap-3 bg-white/10 rounded-xl p-4">
                <Clock className="w-8 h-8 text-yellow-400" />
                <div>
                  <p className="text-purple-200 text-sm">Нужно отработать</p>
                  <p className="text-2xl font-bold text-white">
                    {purchase.formatted}
                  </p>
                  <p className="text-purple-300 text-sm">
                    ({purchase.hours.toFixed(1)} часов)
                  </p>
                </div>
              </div>

              {purchase.comment && (
                <div className="bg-black/20 rounded-xl p-4 border-l-4 border-purple-500">
                  <p className="text-purple-100 italic">{purchase.comment}</p>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default App;
