import { karmaLevelInfo } from '../lib/karma'
import type { KarmaState } from '../lib/karma'

interface KarmaBadgeProps {
  karma: KarmaState
  onClick?: () => void
}

export function KarmaBadge({ karma, onClick }: KarmaBadgeProps) {
  const info = karmaLevelInfo[karma.level]
  const percentage = karma.score

  return (
    <button
      onClick={onClick}
      className="flex items-center gap-2 bg-black/20 hover:bg-black/30 backdrop-blur-sm px-3 py-2 rounded-xl transition-all cursor-pointer"
      title={info.description}
    >
      <span className="text-xl">{info.icon}</span>
      <div className="flex flex-col items-start">
        <span className="text-xs text-purple-200 font-medium">
          {info.displayName}
        </span>
        <div className="flex items-center gap-2">
          <div className="w-20 h-2 bg-gray-700 rounded-full overflow-hidden">
            <div
              className="h-full rounded-full transition-all duration-500"
              style={{
                width: `${percentage}%`,
                backgroundColor: info.color,
              }}
            />
          </div>
          <span className="text-xs text-white font-bold">
            {Math.round(percentage)}
          </span>
        </div>
      </div>
    </button>
  )
}
