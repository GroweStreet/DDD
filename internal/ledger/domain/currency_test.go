package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrency_Unknown(t *testing.T) {

	_, err := NewCurrency("XXX")
	assert.ErrorIs(t, err, ErrUnknownCurrency)
}

func TestCurrency_Known(t *testing.T) {

	usd, err := NewCurrency("USD")
	assert.Nil(t, err)

	assert.Equal(t, "USD", usd.Code())
	assert.Equal(t, int32(2), usd.MinorUnits())
	assert.Equal(t, "$", usd.Symbol())
}

func TestCurrency_JPYNoMinorUnits(t *testing.T) {

	jpy, err := NewCurrency("JPY")
	assert.Nil(t, err)

	assert.Equal(t, int32(0), jpy.MinorUnits())
}
