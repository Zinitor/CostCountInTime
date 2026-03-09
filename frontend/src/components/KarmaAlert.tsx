import { Award, Sparkles } from 'lucide-react'
import { useEffect, useState } from 'react'
import type { KarmaLevel } from '../lib/karma'
import { karmaLevelInfo } from '../lib/karma'

interface KarmaAlertProps {
  level: KarmaLevel
  score: number
  scoreDelta: number
  purchaseComment: string
  onClose: () => void
}

export function KarmaAlert({
  level,
  score,
  scoreDelta,
  purchaseComment,
  onClose,
}: KarmaAlertProps) {
  const [visible, setVisible] = useState(false)
  const info = karmaLevelInfo[level]

  useEffect(() => {
    // Запускаем анимацию появления
    requestAnimationFrame(() => {
      setVisible(true)
    })

    // Автозакрытие через 5 секунд
    const timer = setTimeout(() => {
      handleClose()
    }, 5000)

    return () => clearTimeout(timer)
  }, [])

  const handleClose = () => {
    setVisible(false)
    setTimeout(onClose, 300)
  }

  const isPositive = scoreDelta <= 0
  const isExtreme = score > 80 || score < 20

  return (
    <div
      className={`
        fixed top-4 left-1/2 -translate-x-1/2 z-50
        transition-all duration-500 ease-out
        ${visible ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-4'}
      `}
    >
      <div
        className={`
          relative overflow-hidden
          bg-gradient-to-r ${isPositive ? 'from-green-600 to-emerald-600' : 'from-orange-600 to-red-600'}
          rounded-2xl shadow-2xl p-4 pr-12
          min-w-[320px] max-w-md
          ${isExtreme ? 'animate-karma-glow' : ''}
        `}
      >
        {/* Фоновые искры */}
        <div className="absolute inset-0 overflow-hidden pointer-events-none">
          {[...Array(6)].map((_, i) => (
            <Sparkles
              key={i}
              className="absolute text-white/20 w-4 h-4 animate-pulse"
              style={{
                left: `${15 + i * 15}%`,
                top: `${20 + (i % 3) * 30}%`,
                animationDelay: `${i * 0.2}s`,
              }}
            />
          ))}
        </div>

        <div className="relative z-10">
          <div className="flex items-start gap-3">
            <div className="text-3xl">{info.icon}</div>
            <div className="flex-1">
              <div className="flex items-center gap-2 mb-1">
                <Award className="w-5 h-5 text-white/80" />
                <span className="text-white font-bold text-sm">
                  {info.displayName}
                </span>
              </div>
              <p className="text-white/90 text-sm mb-2">
                {purchaseComment}
              </p>
              <div className="flex items-center gap-3">
                <div className="flex-1 h-2 bg-black/30 rounded-full overflow-hidden">
                  <div
                    className={`h-full rounded-full transition-all duration-1000 ${
                      isPositive ? 'bg-green-300' : 'bg-red-300'
                    }`}
                    style={{ width: `${score}%` }}
                  />
                </div>
                <span className="text-white font-bold text-sm min-w-[3ch]">
                  {Math.round(score)}
                </span>
              </div>
              {scoreDelta !== 0 && (
                <p className={`text-xs mt-1 ${
                  isPositive ? 'text-green-200' : 'text-red-200'
                }`}>
                  {scoreDelta > 0 ? '+' : ''}{scoreDelta.toFixed(1)} к карме
                </p>
              )}
            </div>
          </div>
        </div>

        {/* Кнопка закрытия */}
        <button
          onClick={handleClose}
          className="absolute top-2 right-2 text-white/60 hover:text-white transition-colors p-1"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>
  )
}
