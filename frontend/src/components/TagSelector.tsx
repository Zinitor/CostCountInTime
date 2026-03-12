import { Tag } from "../lib/tags";

interface TagSelectorProps {
  tags: Tag[];
  selectedTag: string | null;
  onSelectTag: (tagId: string | null) => void;
}

export function TagSelector({ tags, selectedTag, onSelectTag }: TagSelectorProps) {
  return (
    <div className="w-full">
      <label className="block text-purple-200 text-sm mb-2">
        Тег <span className="text-purple-400">(опционально)</span>
      </label>
      <div className="flex gap-2 overflow-x-auto pb-2 scrollbar-thin scrollbar-thumb-purple-500/50 scrollbar-track-transparent">
        {/* Кнопка "Без тега" */}
        <button
          type="button"
          onClick={() => onSelectTag(null)}
          className={`
            flex-shrink-0 w-16 h-16 rounded-full border-2 transition-all
            flex items-center justify-center text-xs font-medium
            ${
              selectedTag === null
                ? 'border-white bg-white/20 text-white'
                : 'border-white/30 bg-white/5 text-purple-200 hover:bg-white/10'
            }
          `}
        >
          Нет
        </button>
        
        {/* Кнопки тегов */}
        {tags.map((tag) => (
          <button
            key={tag.id}
            type="button"
            onClick={() => onSelectTag(tag.id)}
            className={`
              flex-shrink-0 w-16 h-16 rounded-full border-2 transition-all
              flex items-center justify-center text-xs font-medium
              ${
                selectedTag === tag.id
                  ? 'border-white scale-110 shadow-lg'
                  : 'border-transparent hover:scale-105'
              }
            `}
            style={{
              background: selectedTag === tag.id
                ? `linear-gradient(135deg, ${tag.colorStart}, ${tag.colorEnd})`
                : `linear-gradient(135deg, ${tag.colorStart}80, ${tag.colorEnd}80)`,
              color: 'white',
            }}
            title={tag.name}
          >
            {tag.name}
          </button>
        ))}
      </div>
    </div>
  );
}
