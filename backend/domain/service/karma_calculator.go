package service

import (
	"math/rand"
	"time"

	"costcountintimeapp/domain/entity"
	"costcountintimeapp/domain/repository"
	"costcountintimeapp/domain/valueobject"
)

// KarmaCalculatorService отвечает за расчёт изменений кармы
type KarmaCalculatorService struct {
	rand    *rand.Rand
	tagRepo repository.TagRepository // Репозиторий тегов
}

// KarmaCalculationResult содержит результат расчёта кармы
type KarmaCalculationResult struct {
	Delta           float64
	Reasons         []string
	PurchaseComment string
}

// NewKarmaCalculatorService создаёт новый сервис расчёта кармы
func NewKarmaCalculatorService(tagRepo repository.TagRepository) *KarmaCalculatorService {
	return &KarmaCalculatorService{
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
		tagRepo: tagRepo,
	}
}

// Комментарии для разных диапазонов цен
var commentsByRange = map[string][]string{
	"cheap": {
		"😏 Неплохо, это всего лишь обед в ресторане!",
		"💪 Да ты богач, даже не заметишь!",
		"⚡ Быстрее, чем сходить в магазин за продуктами!",
		"🎯 Бери не думай, это копейки!",
	},
	"reasonable": {
		"🤔 Нормальная покупка, но подумай дважды.",
		"💸 Неделя работы — не шутка, но терпимо.",
		"📊 Соотношение цена/качество того стоит?",
		"⏳ Пару дней работы — решай сам.",
	},
	"expensive": {
		"😳 Серьёзно? Это же целая зарплата!",
		"🤬 Ты долбоёб за столько брать?",
		"💀 Месяц жизни за эту штуку? Оно того стоит?",
		"🚩 Красный флаг: ты работаешь на вещь, а не она на тебя.",
		"🧐 А точно нельзя найти дешевле?",
	},
	"crazy": {
		"🤯 ТЫ УВЕРЕН ЧТО ТЕБЕ ВООБЩЕ НУЖНА ЭТА ШТУКА?",
		"😱 Несколько месяцев работы! Ты в своём уме?!",
		"💸 Это не покупка, это финансовое самоубийство!",
		"🛑 СТОП! Может лучше отложишь и купишь позже?",
		"🎭 Ты работаешь слишком много, чтобы тратить так бездумно!",
	},
	"insane": {
		"🚨 ГОСПОДИ, ДА ТЫ МАНЬЯК КАКОЙ-ТО!",
		"💀💀💀 Это не покупка — это приговор!",
		"🏦 Может проще банк ограбить?",
		"😭 Бедный ты человек... искренне сочувствую.",
		"🪦 RIP твои нервы и свободное время.",
	},
	"very_insane": {
		"🔥 ФИНАНСОВЫЙ АПОКАЛИПСИС!",
		"💀💀💀 ЭТО НЕ ПОКУПКА — ЭТО ПРИГОВОР!",
		"🤡 КЛОУН ФИНАНСОВЫЙ АКТИВИРОВАЛСЯ!",
		"🚨 ТРЕВОГА! ВЫХОДИТЕ ИЗ МАГАЗИНА!",
		"😭 Я ПЛАЧУ ЗА ТВОЙ КОШЕЛЁК!",
	},
}

// Calculate рассчитывает изменение кармы для покупки
func (s *KarmaCalculatorService) Calculate(price, hoursNeeded float64, itemName string, lastEvent *entity.KarmaEvent, stats KarmaStats, tagID string) KarmaCalculationResult {
	var delta float64
	var reasons []string
	var purchaseComment string

	// 1. Цена относительно зарплаты (часы работы)
	delta, reasons, purchaseComment = s.calculatePriceDelta(hoursNeeded, price)

	// 2. Если выбран тег, используем персонализированный комментарий из тега
	if tagID != "" && s.tagRepo != nil {
		tag, err := s.tagRepo.GetByID(tagID)
		if err == nil && tag != nil {
			// Используем комментарий из тега вместо дефолтного
			purchaseComment = tag.GetRandomComment()
			reasons = append(reasons, "Тег: "+tag.Name)
		}
	}

	// 3. Импульсивность (проверяем последнюю покупку)
	if lastEvent != nil {
		timeSinceLast := time.Now().Unix() - lastEvent.Timestamp
		if timeSinceLast < 3600 { // < 1 часа
			delta += 10
			reasons = append(reasons, "Импульсивная покупка (< 1 часа)")
		} else if timeSinceLast < 86400 { // < 24 часов
			delta += 3
			reasons = append(reasons, "Покупка в тот же день")
		}
	}

	// 4. Накопленная усталость (сумма за 24 часа)
	if stats.SpentToday > 50000 { // > 50к за день
		delta += 20
		reasons = append(reasons, "Много потрачено за день")
	}

	// 5. Паттерны (3+ покупки за день)
	if stats.PurchasesToday >= 3 {
		delta += 5
		reasons = append(reasons, "Много покупок сегодня")
	}

	return KarmaCalculationResult{
		Delta:           delta,
		Reasons:         reasons,
		PurchaseComment: purchaseComment,
	}
}

func (s *KarmaCalculatorService) calculatePriceDelta(hoursNeeded, price float64) (float64, []string, string) {
	var delta float64
	var reasons []string

	rangeKey, _ := s.getRangeByHours(hoursNeeded)
	commentsList := commentsByRange[rangeKey]
	comment := commentsList[s.rand.Intn(len(commentsList))]

	switch rangeKey {
	case "cheap":
		delta -= 2
		reasons = append(reasons, "Разумная трата (< 1 часа)")
	case "reasonable":
		delta += 0
		reasons = append(reasons, "Нормальная покупка (1-8 часов)")
	case "expensive":
		delta += 5
		reasons = append(reasons, "Дороговато (8-40 часов)")
	case "crazy":
		delta += 15
		reasons = append(reasons, "Серьёзная трата (40-200 часов)")
	case "insane":
		delta += 25
		reasons = append(reasons, "Безумие (200-500 часов)")
	case "very_insane":
		delta += 35
		reasons = append(reasons, "Финансовое самоубийство (> 500 часов)")
	}

	return delta, reasons, comment
}

func (s *KarmaCalculatorService) getRangeByHours(hours float64) (string, []string) {
	switch {
	case hours < 2:
		return "cheap", commentsByRange["cheap"]
	case hours < 40:
		return "reasonable", commentsByRange["reasonable"]
	case hours < 200:
		return "expensive", commentsByRange["expensive"]
	case hours < 500:
		return "crazy", commentsByRange["crazy"]
	case hours < 1000:
		return "insane", commentsByRange["insane"]
	default:
		return "very_insane", commentsByRange["very_insane"]
	}
}

// GetLevelByScore возвращает уровень кармы по счёту
func (s *KarmaCalculatorService) GetLevelByScore(score float64) valueobject.LevelInfo {
	for _, level := range valueobject.KarmaLevels {
		if score >= level.MinScore && score < level.MaxScore {
			return level
		}
	}
	return valueobject.KarmaLevels[len(valueobject.KarmaLevels)-1]
}

// GetCommentByLevel возвращает комментарий для уровня кармы
func (s *KarmaCalculatorService) GetCommentByLevel(level valueobject.KarmaLevel, delta float64, reasons []string) string {
	comments := map[valueobject.KarmaLevel][]string{
		valueobject.FinancialGod: {
			"😇 Отличный выбор, разумная трата!",
			"🧘 Ты мастер финансов!",
			"💰 Так держать!",
		},
		valueobject.WiseSpender: {
			"👍 Нормально, но можно ещё лучше",
			"🤔 Подумай, можно ли дешевле",
			"📊 В целом разумно",
		},
		valueobject.NormalPerson: {
			"😐 Бывало и лучше",
			"🤷 Нормально для обычного человека",
			"⚖️ Баланс между разумом и эмоциями",
		},
		valueobject.RiskyBuyer: {
			"😬 Чувак, ты там не перегибаешь?",
			"🚩 Красный флаг! Остановись!",
			"💸 К тебе движутся финансовые проблемы",
		},
		valueobject.FinancialDegen: {
			"🤡 Я вызываю финансового консультанта, у тебя проблемы!",
			"🚨 ТРЕВОГА! Твой кошелёк в опасности!",
			"💀💀💀 Это не покупки, это финансовое самоубийство!",
		},
	}

	levelComments := comments[level]
	if len(levelComments) == 0 {
		return ""
	}

	// Добавляем контекст если была большая дельта
	if delta > 15 {
		reason := ""
		if len(reasons) > 0 {
			reason = reasons[0]
		}
		return levelComments[0] + " (" + reason + ")"
	}

	return levelComments[s.rand.Intn(len(levelComments))]
}

// KarmaStats содержит статистику для расчёта кармы
type KarmaStats struct {
	TotalSpent     float64
	TotalHours     float64
	AvgPurchase    float64
	ImpulseCount   int
	BestStreak     int
	CurrentStreak  int
	PurchasesToday int
	SpentToday     float64
}
