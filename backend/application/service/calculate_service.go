package service

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"costcountintimeapp/application/dto"
	"costcountintimeapp/domain/repository"
	"costcountintimeapp/domain/valueobject"
)

// CalculateService отвечает за расчёт времени работы для покупки
type CalculateService struct {
	rand    *rand.Rand
	tagRepo repository.TagRepository
}

// NewCalculateService создаёт новый CalculateService
func NewCalculateService(tagRepo repository.TagRepository) *CalculateService {
	return &CalculateService{
		rand:    rand.New(rand.NewSource(time.Now().UnixNano())),
		tagRepo: tagRepo,
	}
}

// Calculate рассчитывает время работы для покупки
func (s *CalculateService) Calculate(req dto.CalculateRequest) (*dto.CalculateResponse, error) {
	// Валидация входных данных
	if req.Salary <= 0 || req.Price <= 0 {
		return nil, fmt.Errorf("salary and price must be positive")
	}

	period, err := valueobject.NewTimePeriod(req.SalaryPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid salary period: %w", err)
	}

	// Расчёт
	hoursInPeriod := period.HoursInPeriod()
	salaryPerHour := req.Salary / hoursInPeriod
	hoursNeeded := req.Price / salaryPerHour
	hoursNeeded = math.Round(hoursNeeded*100) / 100
	daysNeeded := math.Round((hoursNeeded/8)*100) / 100

	// Получаем информацию о теге если указан
	var tagName string
	if req.TagID != "" {
		tag, err := s.tagRepo.GetByID(req.TagID)
		if err == nil && tag != nil {
			tagName = tag.Name
		}
	}

	return &dto.CalculateResponse{
		Hours:     hoursNeeded,
		Days:      daysNeeded,
		Formatted: s.formatTime(hoursNeeded),
		ItemName:  req.ItemName,
		Price:     req.Price,
		Comment:   s.getRandomComment(hoursNeeded),
		TagID:     req.TagID,
		TagName:   tagName,
	}, nil
}

func (s *CalculateService) formatTime(hours float64) string {
	if hours < 1 {
		mins := int(math.Round(hours * 60))
		return fmt.Sprintf("%d мин", mins)
	}

	days := int(hours / 8)
	remainingHours := int(math.Round(hours - float64(days)*8))

	if days > 0 && remainingHours > 0 {
		return fmt.Sprintf("%d дн %d ч", days, remainingHours)
	} else if days > 0 {
		return fmt.Sprintf("%d дн", days)
	}
	return fmt.Sprintf("%d ч", int(math.Round(hours)))
}

func (s *CalculateService) getRandomComment(hours float64) string {
	var rangeKey string

	switch {
	case hours < 2:
		rangeKey = "cheap"
	case hours < 40:
		rangeKey = "reasonable"
	case hours < 200:
		rangeKey = "expensive"
	case hours < 500:
		rangeKey = "crazy"
	default:
		rangeKey = "insane"
	}

	commentsByRange := map[string][]string{
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
	}

	comments := commentsByRange[rangeKey]
	if len(comments) == 0 {
		return ""
	}
	return comments[s.rand.Intn(len(comments))]
}
