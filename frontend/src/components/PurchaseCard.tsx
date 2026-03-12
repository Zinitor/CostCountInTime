import { Clock, Trash2 } from "lucide-react";
import { useEffect, useState } from "react";

interface PurchaseCardProps {
  id: number;
  itemName: string;
  price: number;
  hours: number;
  formatted: string;
  comment: string;
  purchaseComment?: string;
  onRemove: (id: number) => void;
  isNew?: boolean;
  highlightKarma?: boolean;
  tagColorStart?: string;
  tagColorEnd?: string;
}

export function PurchaseCard({
  id,
  itemName,
  price,
  hours,
  formatted,
  comment,
  purchaseComment,
  onRemove,
  isNew = false,
  highlightKarma = false,
  tagColorStart,
  tagColorEnd,
}: PurchaseCardProps) {
  const [animationClass, setAnimationClass] = useState("");
  const [removing, setRemoving] = useState(false);

  useEffect(() => {
    if (isNew) {
      setAnimationClass("animate-slide-in-left");
    }
  }, [isNew]);

  useEffect(() => {
    if (highlightKarma) {
      setAnimationClass("animate-karma-shake animate-karma-glow");
      const timer = setTimeout(() => {
        setAnimationClass("");
      }, 2000);
      return () => clearTimeout(timer);
    }
  }, [highlightKarma]);

  const handleRemove = () => {
    setRemoving(true);
    setTimeout(() => {
      onRemove(id);
    }, 300);
  };

  // Определяем градиент в зависимости от тега
  const gradientStyle =
    tagColorStart && tagColorEnd
      ? {
          background: `linear-gradient(135deg, ${tagColorStart}30, ${tagColorEnd}30)`,
        }
      : {};

  const borderColor =
    tagColorStart && tagColorEnd ? { borderColor: `${tagColorStart}60` } : {};

  return (
    <div
      className={`
        backdrop-blur-lg rounded-2xl p-6 border border-white/20
        transition-all duration-300
        ${animationClass}
        ${removing ? "opacity-0 scale-95" : "opacity-100 scale-100"}
        ${highlightKarma ? "ring-2 ring-purple-500 ring-offset-2 ring-offset-transparent" : ""}
      `}
      style={gradientStyle}
    >
      <div className="flex justify-between items-start mb-3">
        <div>
          <h3 className="text-lg font-semibold text-white">{itemName}</h3>
          <p className="text-purple-200">{price.toLocaleString()} ₽</p>
        </div>
        <button
          onClick={handleRemove}
          className="text-purple-300 hover:text-white hover:bg-white/10 p-2 rounded-lg transition-colors"
          title="Удалить"
        >
          <Trash2 className="w-5 h-5" />
        </button>
      </div>

      <div className="flex items-center gap-3 bg-white/10 rounded-xl p-4 mb-3">
        <Clock className="w-8 h-8 text-yellow-400" />
        <div>
          <p className="text-purple-200 text-sm">Нужно отработать</p>
          <p className="text-2xl font-bold text-white">{formatted}</p>
          <p className="text-purple-300 text-sm">({hours.toFixed(1)} часов)</p>
        </div>
      </div>

      {comment && (
        <div
          className="bg-black/20 rounded-xl p-4 border-l-4 border-purple-500 mb-3"
          style={borderColor}
        >
          <p className="text-purple-100 italic">{comment}</p>
        </div>
      )}

      {purchaseComment && (
        <div className="bg-purple-500/10 rounded-xl p-3 border border-purple-500/30">
          <div className="flex items-center gap-2 mb-1">
            <Clock className="w-4 h-4 text-purple-400" />
            <span className="text-purple-300 text-xs font-medium">
              Финансовая карма
            </span>
          </div>
          <p className="text-purple-100 text-sm italic">{purchaseComment}</p>
        </div>
      )}
    </div>
  );
}
