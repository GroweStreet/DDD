package domain

type Currency struct {
	code       string
	minorUnits int32
}

func NewCurrency(code string) (Currency, error) {

	meta, ok := registry[code]
	if !ok {
		return Currency{}, ErrUnknownCurrency
	}

	return Currency{
		code:       code,
		minorUnits: meta.minorUnits,
	}, nil
}

func (c Currency) Code() string {
	return c.code
}

func (c Currency) MinorUnits() int32 {
	return c.minorUnits
}

func (c Currency) Symbol() string {
	return registry[c.code].symbol
}
