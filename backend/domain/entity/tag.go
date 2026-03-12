package entity

// Tag представляет тег для покупки
type Tag struct {
	ID         string
	Name       string
	ColorStart string // Начальный цвет градиента (hex)
	ColorEnd   string // Конечный цвет градиента (hex)
	Comments   []string
}

// Предопределённые теги
var DefaultTags = []Tag{
	{
		ID:         "food",
		Name:       "Еда",
		ColorStart: "#ef4444",
		ColorEnd:   "#dc2626",
		Comments: []string{
			"🥗 Дороговато для еды!",
			"🍽️ Ресторанный чек дома?",
			"🥐 Золотой круассан?",
			"🍣 Суши из золота?",
			"🥩 Мраморная говядина из Японии?",
		},
	},
	{
		ID:         "entertainment",
		Name:       "Развлечения",
		ColorStart: "#f59e0b",
		ColorEnd:   "#d97706",
		Comments: []string{
			"🎮 Игры того стоят!",
			"🎬 Кино с друзьями — отлично!",
			"🎢 Развлечения нужны для жизни!",
			"🎯 Главное — впечатления!",
		},
	},
	{
		ID:         "transport",
		Name:       "Транспорт",
		ColorStart: "#3b82f6",
		ColorEnd:   "#2563eb",
		Comments: []string{
			"🚗 Транспорт — необходимость!",
			"⛽ Бензин дороже с каждым днём!",
			"🚌 Общественный транспорт экономит деньги!",
			"🚕 Такси — иногда можно себе позволить!",
		},
	},
	{
		ID:         "shopping",
		Name:       "Покупки",
		ColorStart: "#ec4899",
		ColorEnd:   "#db2777",
		Comments: []string{
			"🛍️ Шопоголик активировался!",
			"👕 Одежда нужна, но в меру!",
			"💄 Красота требует жертв!",
			"🎁 Подарки — это приятно!",
		},
	},
	{
		ID:         "tech",
		Name:       "Техника",
		ColorStart: "#8b5cf6",
		ColorEnd:   "#7c3aed",
		Comments: []string{
			"💻 Техника — это инвестиция!",
			"📱 Новый гаджет всегда радует!",
			"⌨️ Производительность важна!",
			"🖥️ Технологии будущего!",
		},
	},
}

// NewTag создаёт новый тег
func NewTag(id, name, colorStart, colorEnd string, comments []string) *Tag {
	return &Tag{
		ID:         id,
		Name:       name,
		ColorStart: colorStart,
		ColorEnd:   colorEnd,
		Comments:   comments,
	}
}

// GetRandomComment возвращает случайный комментарий для тега
func (t *Tag) GetRandomComment() string {
	if len(t.Comments) == 0 {
		return ""
	}
	return t.Comments[0] // Для простоты берём первый, можно сделать рандом
}
