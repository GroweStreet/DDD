package domain

type currencyMeta struct {
	minorUnits int32
	symbol     string
}

var registry = map[string]currencyMeta{
	"USD": {minorUnits: 2, symbol: "$"},
	"RUB": {minorUnits: 2, symbol: "₽"},
	"JPY": {minorUnits: 0, symbol: "¥"},
}
