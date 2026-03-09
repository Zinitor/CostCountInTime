package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"

	"costcountintimeapp/models"
)

var HoursInPeriod = map[string]float64{
	"hour":  1,
	"day":   8,
	"month": 168,
}

// Комментарии для разных диапазонов времени
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
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CalculateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Salary <= 0 || req.Price <= 0 {
		http.Error(w, "Salary and price must be positive", http.StatusBadRequest)
		return
	}

	hoursInPeriod, ok := HoursInPeriod[req.SalaryPeriod]
	if !ok {
		http.Error(w, "Invalid salary period", http.StatusBadRequest)
		return
	}

	salaryPerHour := req.Salary / hoursInPeriod
	hoursNeeded := req.Price / salaryPerHour
	hoursNeeded = math.Round(hoursNeeded*100) / 100
	daysNeeded := math.Round((hoursNeeded/8)*100) / 100

	response := models.CalculateResponse{
		Hours:     hoursNeeded,
		Days:      daysNeeded,
		Formatted: formatTime(hoursNeeded),
		ItemName:  req.ItemName,
		Price:     req.Price,
		Comment:   getRandomComment(hoursNeeded),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getRandomComment(hours float64) string {
	var rangeKey string

	switch {
	case hours < 2:
		rangeKey = "cheap" // < 2 часов
	case hours < 40:
		rangeKey = "reasonable" // < 5 дней
	case hours < 200:
		rangeKey = "expensive" // < 25 дней (~1 месяц)
	case hours < 500:
		rangeKey = "crazy" // < 62 дней (~2-3 месяца)
	default:
		rangeKey = "insane" // > 2-3 месяцев
	}

	comments := commentsByRange[rangeKey]
	if len(comments) == 0 {
		return ""
	}
	return comments[rand.Intn(len(comments))]
}

func formatTime(hours float64) string {
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

func StaticHandler(fs http.FileSystem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		filePath := filepath.Clean("frontend/dist" + path)
		f, err := fs.Open(filePath)
		if err != nil {
			http.ServeFile(w, r, "frontend/dist/index.html")
			return
		}
		defer f.Close()
		http.ServeFile(w, r, filePath)
	}
}

func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}
