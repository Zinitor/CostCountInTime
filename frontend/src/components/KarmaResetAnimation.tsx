import { useEffect, useState } from "react";
import { Sparkles, Star, Zap } from "lucide-react";

interface KarmaResetAnimationProps {
  onComplete: () => void;
}

interface Particle {
  id: number;
  left: number;
  delay: number;
  size: number;
  rotation: number;
}

interface FallingStar {
  id: number;
  left: number;
  delay: number;
  size: number;
}

export function KarmaResetAnimation({ onComplete }: KarmaResetAnimationProps) {
  const [phase, setPhase] = useState<
    "darkness" | "shake" | "flash" | "explosion" | "purify" | "stars" | "done"
  >("darkness");
  const [particles, setParticles] = useState<Particle[]>([]);
  const [fallingStars, setFallingStars] = useState<FallingStar[]>([]);

  useEffect(() => {
    // Генерируем частицы
    const newParticles: Particle[] = [];
    for (let i = 0; i < 20; i++) {
      newParticles.push({
        id: i,
        left: Math.random() * 100,
        delay: Math.random() * 0.5,
        size: 4 + Math.random() * 8,
        rotation: Math.random() * 360,
      });
    }
    setParticles(newParticles);

    // Генерируем падающие звёзды
    const newStars: FallingStar[] = [];
    for (let i = 0; i < 15; i++) {
      newStars.push({
        id: i,
        left: Math.random() * 100,
        delay: Math.random() * 0.8,
        size: 8 + Math.random() * 12,
      });
    }
    setFallingStars(newStars);

    // Фазы анимации
    const shakeTimer = setTimeout(() => setPhase("shake"), 100);
    const flashTimer = setTimeout(() => setPhase("flash"), 600);
    const explosionTimer = setTimeout(() => setPhase("explosion"), 800);
    const purifyTimer = setTimeout(() => setPhase("purify"), 1500);
    const starsTimer = setTimeout(() => setPhase("stars"), 2500);
    const doneTimer = setTimeout(() => {
      setPhase("done");
      onComplete();
    }, 4000);

    return () => {
      clearTimeout(shakeTimer);
      clearTimeout(flashTimer);
      clearTimeout(explosionTimer);
      clearTimeout(purifyTimer);
      clearTimeout(starsTimer);
      clearTimeout(doneTimer);
    };
  }, [onComplete]);

  return (
    <div className="fixed inset-0 z-[100] pointer-events-none overflow-hidden">
      {/* Фаза 1: Тьма исчезает */}
      {phase === "darkness" && (
        <div className="absolute inset-0 bg-black animate-reset-darkness-fade" />
      )}

      {/* Фаза 2: Тряска экрана */}
      {phase === "shake" && (
        <div className="absolute inset-0 animate-reset-shake">
          <div className="absolute inset-0 bg-gradient-to-r from-red-900/50 via-orange-900/50 to-red-900/50" />
        </div>
      )}

      {/* Фаза 3: Вспышка света */}
      {phase === "flash" && (
        <>
          <div className="absolute inset-0 bg-white animate-reset-flash" />
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="w-96 h-96 bg-white rounded-full blur-3xl animate-reset-flash" />
          </div>
        </>
      )}

      {/* Фаза 4: Взрывная волна */}
      {phase === "explosion" && (
        <div className="absolute inset-0 flex items-center justify-center">
          <div className="w-32 h-32 bg-gradient-to-r from-yellow-400 via-orange-500 to-red-500 rounded-full blur-xl animate-reset-explosion" />
          <div
            className="absolute w-24 h-24 bg-white rounded-full blur-lg animate-reset-explosion"
            style={{ animationDelay: "0.1s" }}
          />
        </div>
      )}

      {/* Фаза 5: Вращающиеся частицы */}
      {(phase === "explosion" || phase === "purify") && (
        <div className="absolute inset-0 flex items-center justify-center">
          {particles.map((particle) => (
            <div
              key={particle.id}
              className="absolute"
              style={{
                left: "50%",
                top: "50%",
                width: particle.size,
                height: particle.size,
                backgroundColor: [
                  "#22c55e",
                  "#84cc16",
                  "#eab308",
                  "#f97316",
                  "#ef4444",
                ][Math.floor(Math.random() * 5)],
                borderRadius: Math.random() > 0.5 ? "50%" : "0",
                transform: `rotate(${particle.rotation}deg)`,
                animation:
                  "reset-particles 1.2s cubic-bezier(0.19, 1, 0.22, 1) forwards",
                animationDelay: `${particle.delay}s`,
                opacity: 0.8,
              }}
            />
          ))}
        </div>
      )}

      {/* Фаза 6: Пульсация очищения */}
      {phase === "purify" && (
        <div className="absolute inset-0 flex items-center justify-center animate-reset-purify">
          <div className="w-[600px] h-[600px] bg-gradient-to-r from-green-400 via-emerald-500 to-teal-500 rounded-full blur-3xl opacity-50" />
          <div className="absolute text-center">
            <Sparkles className="w-24 h-24 text-green-300 mx-auto mb-4 animate-pulse" />
            <div className="text-4xl font-bold text-green-200 mb-2">
              ✨ КАРМА ОЧИЩЕНА ✨
            </div>
            <div className="text-green-300 text-lg">
              Начинаем с чистого листа!
            </div>
          </div>
        </div>
      )}

      {/* Фаза 7: Падающие звёзды */}
      {(phase === "stars" || phase === "done") && (
        <div className="absolute inset-0">
          {fallingStars.map((star) => (
            <div
              key={star.id}
              className="absolute"
              style={{
                left: `${star.left}%`,
                top: "-50px",
                animation:
                  "reset-falling-star 1.5s cubic-bezier(0.19, 1, 0.22, 1) forwards",
                animationDelay: `${star.delay}s`,
              }}
            >
              <Star
                className="text-yellow-300 fill-yellow-300"
                style={{
                  width: star.size,
                  height: star.size,
                }}
              />
            </div>
          ))}

          {/* Дополнительные искры */}
          {[...Array(30)].map((_, i) => (
            <div
              key={i}
              className="absolute"
              style={{
                left: `${Math.random() * 100}%`,
                top: `${Math.random() * 100}%`,
                animation:
                  "reset-spiral 1.5s cubic-bezier(0.19, 1, 0.22, 1) forwards",
                animationDelay: `${Math.random() * 0.5}s`,
              }}
            >
              <Zap
                className="text-yellow-400 fill-yellow-400"
                style={{
                  width: 12 + Math.random() * 16,
                  height: 12 + Math.random() * 16,
                }}
              />
            </div>
          ))}

          {/* Финальное сообщение */}
          <div className="absolute inset-0 flex items-center justify-center">
            <div className="text-center animate-fade-in-down">
              <div className="text-6xl mb-4">😇</div>
              <div className="text-3xl font-bold text-green-200 mb-2">
                Финансовый Бог!
              </div>
              <div className="text-green-300">Карма сброшена до 0</div>
            </div>
          </div>
        </div>
      )}

      {/* Восходящий свет */}
      {phase === "purify" && (
        <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-green-500/30 to-transparent animate-reset-rising-light" />
      )}
    </div>
  );
}
