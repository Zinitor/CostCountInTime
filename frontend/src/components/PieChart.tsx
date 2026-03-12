import { useEffect, useState } from "react";

interface PieChartProps {
  totalBudget: number;
  spentAmount: number;
  salaryPerHour: number;
  segments: PieSegment[];
}

interface PieSegment {
  value: number;
  color: string;
  label: string;
}

interface AnimatedSegment extends PieSegment {
  delay: number;
  startAngle: number;
  endAngle: number;
}

export function PieChart({
  totalBudget,
  spentAmount,
  salaryPerHour,
  segments,
}: PieChartProps) {
  const [animatedSegments, setAnimatedSegments] = useState<AnimatedSegment[]>(
    [],
  );
  const remainingAmount = Math.max(0, totalBudget - spentAmount);
  const savedHours = remainingAmount / salaryPerHour;

  // Форматируем часы в читаемый формат
  const formatHours = (hours: number): string => {
    if (hours < 1) {
      const mins = Math.round(hours * 60);
      return `${mins} мин`;
    }
    const days = Math.floor(hours / 8);
    const remainingHours = Math.round(hours % 8);
    if (days > 0 && remainingHours > 0) {
      return `${days} дн ${remainingHours} ч`;
    } else if (days > 0) {
      return `${days} дн`;
    }
    return `${Math.round(hours)} ч`;
  };

  // Комментарий в зависимости от процента сэкономленного
  const savedPercent =
    totalBudget > 0 ? (remainingAmount / totalBudget) * 100 : 0;
  const getComment = () => {
    if (savedPercent >= 80) {
      return "💰 Отлично! Ты очень экономный!";
    } else if (savedPercent >= 50) {
      return "👍 Хороший баланс трат и сбережений!";
    } else if (savedPercent >= 20) {
      return "⚠️ Осторожно, траты растут!";
    } else if (savedPercent > 0) {
      return "🚨 Почти всё потрачено!";
    } else {
      return "💸 Всё потрачено!";
    }
  };

  // Рассчитываем сегменты при изменении данных
  useEffect(() => {
    const pieSegments: AnimatedSegment[] = [];
    let currentAngle = 0;

    // Защита от division by zero
    if (totalBudget <= 0) {
      setAnimatedSegments([]);
      return;
    }

    // Добавляем потраченные сегменты
    segments.forEach((segment, index) => {
      if (segment.value > 0) {
        const angle = (segment.value / totalBudget) * 360;
        pieSegments.push({
          ...segment,
          delay: index * 100,
          startAngle: currentAngle,
          endAngle: currentAngle + angle,
        });
        currentAngle += angle;
      }
    });

    // Добавляем оставшиеся средства (зелёный сектор)
    if (remainingAmount > 0) {
      const angle = (remainingAmount / totalBudget) * 360;
      pieSegments.push({
        value: remainingAmount,
        color: "#22c55e",
        label: "Остаток",
        delay: segments.length * 100,
        startAngle: currentAngle,
        endAngle: currentAngle + angle,
      });
    }

    setAnimatedSegments(pieSegments);
  }, [segments, remainingAmount, totalBudget]);

  const polarToCartesian = (
    cx: number,
    cy: number,
    r: number,
    angleInDegrees: number,
  ): { x: number; y: number } => {
    const angleInRadians = ((angleInDegrees - 90) * Math.PI) / 180.0;
    return {
      x: cx + r * Math.cos(angleInRadians),
      y: cy + r * Math.sin(angleInRadians),
    };
  };

  const createArcPath = (
    cx: number,
    cy: number,
    r: number,
    startAngle: number,
    endAngle: number,
  ): string => {
    // Если угол 360 градусов, рисуем полный круг
    if (endAngle - startAngle >= 359.9) {
      return `M ${cx + r} ${cy} A ${r} ${r} 0 1 0 ${cx - r} ${cy} A ${r} ${r} 0 1 0 ${cx + r} ${cy}`;
    }

    const start = polarToCartesian(cx, cy, r, endAngle);
    const end = polarToCartesian(cx, cy, r, startAngle);
    const largeArcFlag = endAngle - startAngle <= 180 ? "0" : "1";

    return [
      "M",
      cx,
      cy,
      "L",
      start.x,
      start.y,
      "A",
      r,
      r,
      0,
      largeArcFlag,
      0,
      end.x,
      end.y,
      "Z",
    ].join(" ");
  };

  return (
    <div className="bg-white/10 backdrop-blur-lg rounded-2xl p-6 border border-white/20">
      <h3 className="text-xl font-semibold text-white mb-4 text-center">
        Распределение бюджета
      </h3>

      <div className="flex flex-col items-center">
        {/* SVG диаграмма */}
        <div className="relative w-48 h-48">
          <svg viewBox="0 0 200 200" className="w-full h-full">
            {animatedSegments.map((segment, index) => {
              const path = createArcPath(
                100,
                100,
                80,
                segment.startAngle,
                segment.endAngle,
              );
              return (
                <path
                  key={index}
                  d={path}
                  fill={segment.color}
                  style={{
                    opacity: 0,
                    animation: `fadeIn 0.5s ease-out ${segment.delay}ms forwards`,
                    transformOrigin: "center",
                  }}
                />
              );
            })}
            {/* Белый круг в центре для создания кольца */}
            <circle cx="100" cy="100" r="50" fill="rgba(0,0,0,0.3)" />
          </svg>

          {/* Центр с оставшейся суммой */}
          <div className="absolute inset-0 flex flex-col items-center justify-center text-white">
            <div className="text-2xl font-bold">
              {remainingAmount.toLocaleString()} ₽
            </div>
            <div className="text-xs text-purple-200 mt-1">
              из {totalBudget.toLocaleString()} ₽
            </div>
          </div>
        </div>

        {/* Легенда */}
        <div className="mt-4 space-y-2 w-full">
          {segments.map(
            (segment, index) =>
              segment.value > 0 && (
                <div
                  key={index}
                  className="flex items-center justify-between text-sm"
                >
                  <div className="flex items-center gap-2">
                    <div
                      className="w-3 h-3 rounded-full"
                      style={{ backgroundColor: segment.color }}
                    />
                    <span className="text-purple-200">{segment.label}</span>
                  </div>
                  <span className="text-white font-medium">
                    {segment.value.toLocaleString()} ₽
                  </span>
                </div>
              ),
          )}
          {remainingAmount > 0 && (
            <div className="flex items-center justify-between text-sm">
              <div className="flex items-center gap-2">
                <div
                  className="w-3 h-3 rounded-full"
                  style={{ backgroundColor: "#22c55e" }}
                />
                <span className="text-purple-200">Остаток</span>
              </div>
              <span className="text-white font-medium">
                {remainingAmount.toLocaleString()} ₽
              </span>
            </div>
          )}
        </div>

        {/* Комментарий */}
        <div className="mt-4 text-center">
          <div className="text-purple-200 text-sm mb-1">Сэкономлено</div>
          <div className="text-white font-semibold text-lg">
            {remainingAmount.toLocaleString()} ₽
          </div>
          <div className="text-purple-300 text-sm">
            ({formatHours(savedHours)})
          </div>
          <div className="mt-2 text-purple-100 italic text-sm">
            {getComment()}
          </div>
        </div>
      </div>

      {/* Стили для анимации */}
      <style>{`
        @keyframes fadeIn {
          from {
            opacity: 0;
          }
          to {
            opacity: 1;
          }
        }
      `}</style>
    </div>
  );
}
