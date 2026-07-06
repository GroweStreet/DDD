package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoney_AddSameCurrency(t *testing.T) {

	usd, err := NewCurrency("USD")
	assert.Nil(t, err)

	sum, err := NewMoney(100, usd).Add(NewMoney(200, usd))
	assert.Nil(t, err)
	assert.Equal(t, int64(300), sum.Amount())
}

func TestMoney_AddMismatch(t *testing.T) {

	usd, err := NewCurrency("USD")
	assert.Nil(t, err)
	rub, err := NewCurrency("RUB")
	assert.Nil(t, err)

	_, err = NewMoney(100, usd).Add(NewMoney(200, rub))
	assert.ErrorIs(t, err, ErrMismatchingMoneyTypes)
}

func TestMoney_SubSameCurrency(t *testing.T) {

	usd, err := NewCurrency("USD")
	assert.Nil(t, err)

	res, err := NewMoney(300, usd).Sub(NewMoney(100, usd))
	assert.Nil(t, err)
	assert.Equal(t, int64(200), res.Amount())
}

func TestMoney_SubMismatch(t *testing.T) {

	usd, err := NewCurrency("USD")
	assert.Nil(t, err)
	rub, err := NewCurrency("RUB")
	assert.Nil(t, err)

	_, err = NewMoney(300, usd).Sub(NewMoney(100, rub))
	assert.ErrorIs(t, err, ErrMismatchingMoneyTypes)
}

func TestMoney_Negate(t *testing.T) {

	usd, err := NewCurrency("USD")
	assert.Nil(t, err)

	assert.Equal(t, int64(-100), NewMoney(100, usd).Negate().Amount())
	assert.Equal(t, int64(100), NewMoney(-100, usd).Negate().Amount())
}

func TestMoney_IsZero(t *testing.T) {

	usd, err := NewCurrency("USD")
	assert.Nil(t, err)

	assert.True(t, Zero(usd).IsZero())
	assert.False(t, NewMoney(1, usd).IsZero())
}

func TestMoney_IsNegative(t *testing.T) {

	usd, err := NewCurrency("USD")
	assert.Nil(t, err)

	assert.True(t, NewMoney(-1, usd).IsNegative())
	assert.False(t, NewMoney(0, usd).IsNegative())
	assert.False(t, NewMoney(1, usd).IsNegative())
}
