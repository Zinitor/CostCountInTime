package valueobject

// KarmaLevel представляет уровень кармы
type KarmaLevel string

const (
	FinancialGod   KarmaLevel = "financial_god"
	WiseSpender    KarmaLevel = "wise_spender"
	NormalPerson   KarmaLevel = "normal_person"
	RiskyBuyer     KarmaLevel = "risky_buyer"
	FinancialDegen KarmaLevel = "financial_degen"
)

// LevelInfo содержит информацию об уровне кармы
type LevelInfo struct {
	Level       KarmaLevel
	MinScore    float64
	MaxScore    float64
	DisplayName string
	Icon        string
	Color       string
}

// KarmaLevels содержит все уровни кармы
var KarmaLevels = []LevelInfo{
	{FinancialGod, 0, 20, "Financial God", "😇", "#22c55e"},
	{WiseSpender, 20, 40, "Wise Spender", "🧘", "#84cc16"},
	{NormalPerson, 40, 60, "Normal Person", "😐", "#eab308"},
	{RiskyBuyer, 60, 80, "Risky Buyer", "😬", "#f97316"},
	{FinancialDegen, 80, 100, "Financial Degen", "🤡", "#ef4444"},
}
