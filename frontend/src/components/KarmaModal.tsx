import { X, RotateCcw } from "lucide-react";
import { karmaLevelInfo } from "../lib/karma";
import type { KarmaState } from "../lib/karma";

interface KarmaModalProps {
  karma: KarmaState;
  onClose: () => void;
  onReset: () => void;
}

export function KarmaModal({ karma, onClose, onReset }: KarmaModalProps) {
  const info = karmaLevelInfo[karma.level];
  const stats = karma.stats;

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div className="bg-gradient-to-br from-indigo-900 via-purple-900 to-pink-800 rounded-2xl p-6 max-w-md w-full border border-white/20 shadow-2xl">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-white flex items-center gap-2">
            {info.icon} Финансовая Карма
          </h2>
          <div className="flex items-center gap-2">
            <button
              onClick={onReset}
              className="text-orange-300 hover:text-orange-200 transition-colors p-2 rounded-lg hover:bg-white/10"
              title="Сбросить карму"
            >
              <RotateCcw className="w-5 h-5" />
            </button>
            <button
              onClick={onClose}
              className="text-purple-300 hover:text-white transition-colors p-2 rounded-lg hover:bg-white/10"
            >
              <X className="w-6 h-6" />
            </button>
          </div>
        </div>

        {/* Score */}
        <div className="bg-black/20 rounded-xl p-4 mb-4">
          <div className="text-center mb-3">
            <div className="text-4xl mb-2">{info.icon}</div>
            <div className="text-xl font-bold text-white">
              {info.displayName}
            </div>
            <div className="text-purple-300 text-sm">{info.description}</div>
          </div>

          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-purple-200">Уровень кармы</span>
              <span className="text-white font-bold">
                {Math.round(karma.score)} / 100
              </span>
            </div>
            <div className="h-4 bg-gray-700 rounded-full overflow-hidden">
              <div
                className="h-full rounded-full transition-all duration-500"
                style={{
                  width: `${karma.score}%`,
                  backgroundColor: info.color,
                }}
              />
            </div>
          </div>
        </div>

        {/* Статистика */}
        <div className="grid grid-cols-2 gap-3">
          <div className="bg-black/20 rounded-xl p-3">
            <div className="text-purple-300 text-xs mb-1">
              Потрачено сегодня
            </div>
            <div className="text-white font-bold text-lg">
              {stats.spent_today.toLocaleString()} ₽
            </div>
          </div>
          <div className="bg-black/20 rounded-xl p-3">
            <div className="text-purple-300 text-xs mb-1">Покупок сегодня</div>
            <div className="text-white font-bold text-lg">
              {stats.purchases_today}
            </div>
          </div>
          <div className="bg-black/20 rounded-xl p-3">
            <div className="text-purple-300 text-xs mb-1">Всего потрачено</div>
            <div className="text-white font-bold text-lg">
              {stats.total_spent.toLocaleString()} ₽
            </div>
          </div>
          <div className="bg-black/20 rounded-xl p-3">
            <div className="text-purple-300 text-xs mb-1">Средняя покупка</div>
            <div className="text-white font-bold text-lg">
              {Math.round(stats.avg_purchase).toLocaleString()} ₽
            </div>
          </div>
        </div>

        {/* Советы */}
        {karma.score > 60 && (
          <div className="mt-4 bg-red-500/20 border border-red-500/50 rounded-xl p-3">
            <div className="text-red-300 text-sm font-medium">
              ⚠️ Зона риска
            </div>
            <div className="text-purple-200 text-xs mt-1">
              Твои траты становятся опасными. Подумай о том, чтобы замедлиться.
            </div>
          </div>
        )}

        {karma.score > 80 && (
          <div className="mt-4 bg-orange-500/20 border border-orange-500/50 rounded-xl p-3">
            <div className="text-orange-300 text-sm font-medium">
              🚨 ОПАСНОСТЬ
            </div>
            <div className="text-purple-200 text-xs mt-1">
              ТВОЙ ФИНАНСОВЫЙ КОНСУЛЬТАНТ ПЛАЧЕТ. СРОЧНО ОСТАНОВИСЬ!
            </div>
          </div>
        )}

        {/* Кнопка сброса */}
        <div className="mt-6 pt-4 border-t border-white/10">
          <button
            onClick={onReset}
            className="w-full bg-orange-500/20 hover:bg-orange-500/30 border border-orange-500/50 text-orange-300 font-medium py-3 px-4 rounded-xl transition-all flex items-center justify-center gap-2"
          >
            <RotateCcw className="w-5 h-5" />
            Сбросить карму
          </button>
          <p className="text-purple-400 text-xs text-center mt-2">
            Удалит всю историю покупок и обнулит счёт
          </p>
        </div>
      </div>
    </div>
  );
}
