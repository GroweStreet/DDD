package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccount_Overdraft(t *testing.T) {

	someTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	currency, err := NewCurrency("USD")
	assert.Nil(t, err)

	money := NewMoney(int64(-6000), currency)

	acc1, err := OpenUserAccount("acc-1", "u-1", currency, someTime)
	assert.Nil(t, err)

	posting, err := NewPosting(acc1.id, money)

	err = acc1.Apply(posting)
	assert.ErrorIs(t, err, ErrOverdraftLimit)
}

func TestAccount_MismatchCurrency(t *testing.T) {

	someTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	currencyUsd, err := NewCurrency("USD")
	assert.Nil(t, err)
	currencyRub, err := NewCurrency("RUB")
	assert.Nil(t, err)

	money := NewMoney(int64(-6000), currencyRub)

	acc1, err := OpenUserAccount("acc-1", "u-1", currencyUsd, someTime)
	assert.Nil(t, err)

	posting, err := NewPosting(acc1.id, money)

	err = acc1.Apply(posting)
	assert.ErrorIs(t, err, ErrDifferentCurrency)
}
