package domain

type Money struct {
	amount   int64
	currency Currency
}

func NewMoney(amount int64, currency Currency) Money {
	return Money{
		amount:   amount,
		currency: currency,
	}
}

func Zero(currency Currency) Money {
	return Money{
		amount:   0,
		currency: currency,
	}
}

func (m Money) Add(other Money) (Money, error) {

	if m.currency != other.currency {
		return Money{}, ErrMismatchingMoneyTypes
	}

	return Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

func (m Money) Sub(other Money) (Money, error) {

	if m.currency != other.currency {
		return Money{}, ErrMismatchingMoneyTypes
	}

	m.amount -= other.Amount()
	return m, nil
}

func (m Money) Negate() Money {

	m.amount = -m.amount
	return m
}

func (m Money) IsZero() bool {
	return m.amount == 0
}

func (m Money) IsNegative() bool {
	return m.amount < 0
}

func (m Money) Amount() int64 {
	return m.amount
}

func (m Money) Currency() Currency {
	return m.currency
}
