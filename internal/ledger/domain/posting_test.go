package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPosting_Zero(t *testing.T) {

	usd, err := NewCurrency("USD")
	assert.Nil(t, err)

	_, err = NewPosting("acc-1", Zero(usd))
	assert.ErrorIs(t, err, ErrZeroPosting)
}

func TestPosting_Accessors(t *testing.T) {

	usd, err := NewCurrency("USD")
	assert.Nil(t, err)

	money := NewMoney(100, usd)

	p, err := NewPosting("acc-1", money)
	assert.Nil(t, err)

	assert.Equal(t, AccountID("acc-1"), p.AccountID())
	assert.Equal(t, money, p.Money())
}
