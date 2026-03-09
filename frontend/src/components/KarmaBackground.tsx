import { useEffect, useState } from 'react'
import type { KarmaLevel } from '../lib/karma'

interface Particle {
  id: number
  emoji: string
  left: number
  delay: number
  duration: number
  size: number
}

interface KarmaBackgroundProps {
  karmaLevel?: KarmaLevel | null
}

const levelEmojis: Record<string, string[]> = {
  financial_god: ['😇', '✨', '💎', '🌟', '💰', '🦋'],
  wise_spender: ['🧘', '🌿', '🍃', '💚', '🌸', '✨'],
  normal_person: ['😐', '💭', '⚖️', '🤔', '💼', '📊'],
  risky_buyer: ['😬', '⚠️', '🔥', '💸', '📉', '😰'],
  financial_degen: ['🤡', '🔥', '💀', '🚨', '😱', '💩'],
}

const levelColors: Record<string, string> = {
  financial_god: 'from-emerald-900 via-green-800 to-teal-900',
  wise_spender: 'from-green-900 via-emerald-800 to-cyan-900',
  normal_person: 'from-indigo-900 via-purple-800 to-pink-900',
  risky_buyer: 'from-orange-900 via-red-800 to-rose-900',
  financial_degen: 'from-red-950 via-orange-900 to-yellow-900',
}

export function KarmaBackground({ karmaLevel }: KarmaBackgroundProps) {
  const [particles, setParticles] = useState<Particle[]>([])
  const level = karmaLevel || 'normal_person'
  const emojis = levelEmojis[level] || levelEmojis.normal_person
  const bgGradient = levelColors[level] || levelColors.normal_person

  useEffect(() => {
    const generateParticles = () => {
      const newParticles: Particle[] = []
      const particleCount = level === 'financial_god' || level === 'wise_spender' ? 20 :
                           level === 'normal_person' ? 15 : 30

      for (let i = 0; i < particleCount; i++) {
        newParticles.push({
          id: i,
          emoji: emojis[Math.floor(Math.random() * emojis.length)],
          left: Math.random() * 100,
          delay: Math.random() * 5,
          duration: 8 + Math.random() * 7,
          size: 12 + Math.random() * 20,
        })
      }
      setParticles(newParticles)
    }

    generateParticles()
  }, [level])

  return (
    <>
      {/* Базовый цвет фона (без анимации, чтобы не моргало) */}
      <div className={`fixed inset-0 bg-gradient-to-br ${bgGradient} -z-20`} />
      
      {/* Анимированный фон */}
      <div className="fixed inset-0 opacity-20 -z-10 pointer-events-none">
        {level === 'financial_god' && <HeavenBackground />}
        {level === 'wise_spender' && <CalmBackground />}
        {level === 'normal_person' && <NeutralBackground />}
        {level === 'risky_buyer' && <WarningBackground />}
        {level === 'financial_degen' && <HellBackground />}
      </div>

      {/* Падающие частицы */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none -z-10">
        {particles.map((particle) => (
          <div
            key={particle.id}
            className="absolute animate-fall"
            style={{
              left: `${particle.left}%`,
              animationDelay: `${particle.delay}s`,
              animationDuration: `${particle.duration}s`,
              fontSize: `${particle.size}px`,
              top: '-50px',
            }}
          >
            {particle.emoji}
          </div>
        ))}
      </div>

      {/* Виньетка */}
      <div className="fixed inset-0 bg-gradient-to-t from-black/40 via-transparent to-black/20 pointer-events-none -z-0" />
    </>
  )
}

// Райский фон
function HeavenBackground() {
  return (
    <div className="absolute inset-0">
      {/* Светящиеся облака */}
      <div className="absolute top-10 left-10 w-32 h-32 bg-white/10 rounded-full blur-3xl animate-float" />
      <div className="absolute top-20 right-20 w-40 h-40 bg-yellow-200/10 rounded-full blur-3xl animate-float-slow" />
      <div className="absolute bottom-1/3 left-1/4 w-24 h-24 bg-emerald-200/10 rounded-full blur-2xl animate-pulse" />
      
      {/* Лучи света */}
      <div className="absolute top-0 left-1/4 w-2 h-full bg-gradient-to-b from-yellow-100/20 to-transparent animate-shimmer" />
      <div className="absolute top-0 right-1/3 w-3 h-full bg-gradient-to-b from-white/10 to-transparent animate-shimmer-delayed" />
    </div>
  )
}

// Спокойный фон
function CalmBackground() {
  return (
    <div className="absolute inset-0">
      <div className="absolute top-1/4 left-1/3 w-48 h-48 bg-green-300/10 rounded-full blur-3xl animate-breathe" />
      <div className="absolute bottom-1/4 right-1/4 w-32 h-32 bg-cyan-300/10 rounded-full blur-2xl animate-breathe-delayed" />
      <div className="absolute top-1/2 left-1/2 w-24 h-24 bg-emerald-200/10 rounded-full blur-xl animate-pulse-slow" />
    </div>
  )
}

// Нейтральный фон
function NeutralBackground() {
  return (
    <div className="absolute inset-0">
      <div className="absolute top-1/3 left-1/4 w-40 h-40 bg-purple-300/10 rounded-full blur-3xl animate-drift" />
      <div className="absolute bottom-1/3 right-1/3 w-36 h-36 bg-pink-300/10 rounded-full blur-2xl animate-drift-slow" />
    </div>
  )
}

// Предупреждающий фон
function WarningBackground() {
  return (
    <div className="absolute inset-0">
      {/* Мерцающие огоньки */}
      <div className="absolute top-1/4 right-1/4 w-32 h-32 bg-orange-500/20 rounded-full blur-2xl animate-flicker" />
      <div className="absolute bottom-1/3 left-1/3 w-24 h-24 bg-red-500/20 rounded-full blur-xl animate-flicker-fast" />
      <div className="absolute top-1/2 right-1/3 w-20 h-20 bg-yellow-500/15 rounded-full blur-lg animate-pulse" />
    </div>
  )
}

// Адский фон
function HellBackground() {
  return (
    <div className="absolute inset-0">
      {/* Огонь */}
      <div className="absolute bottom-0 left-0 right-0 h-1/2 bg-gradient-to-t from-orange-600/30 via-red-600/20 to-transparent animate-fire" />
      
      {/* Пульсирующие огни */}
      <div className="absolute bottom-1/4 left-1/4 w-40 h-40 bg-orange-500/30 rounded-full blur-3xl animate-fire-pulse" />
      <div className="absolute bottom-1/3 right-1/4 w-32 h-32 bg-red-500/30 rounded-full blur-2xl animate-fire-pulse-delayed" />
      <div className="absolute bottom-1/2 left-1/3 w-24 h-24 bg-yellow-500/25 rounded-full blur-xl animate-fire-flicker" />
      
      {/* Искры */}
      <div className="absolute inset-0 overflow-hidden">
        {[...Array(20)].map((_, i) => (
          <div
            key={i}
            className="absolute w-1 h-1 bg-orange-300 rounded-full animate-spark"
            style={{
              left: `${Math.random() * 100}%`,
              top: `${Math.random() * 100}%`,
              animationDelay: `${Math.random() * 3}s`,
            }}
          />
        ))}
      </div>
    </div>
  )
}
