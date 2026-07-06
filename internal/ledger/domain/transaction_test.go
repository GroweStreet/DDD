package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const someAmount = 500
const someAccountID = "acc-1"
const someTransactionID = "t-1"

func TestTransaction_NotEnoughPostings(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currency, err := NewCurrency("USD")
	assert.Nil(t, err)

	money := NewMoney(someAmount, currency)

	posting, err := NewPosting(someAccountID, money)
	assert.Nil(t, err)

	_, err = NewTransaction("t-1", []Posting{posting}, someTime, "key-1", "")
	assert.ErrorIs(t, err, ErrNotEnoughPostings)
}

func TestTransaction_DifferentCurrency(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currencyUsd, err := NewCurrency("USD")
	assert.Nil(t, err)

	currencyRub, err := NewCurrency("RUB")
	assert.Nil(t, err)

	moneyUsd := NewMoney(someAmount, currencyUsd)
	moneyRub := NewMoney(someAmount, currencyRub)

	postingUsd, err := NewPosting(someAccountID, moneyUsd)
	assert.Nil(t, err)

	postingRub, err := NewPosting(someAccountID, moneyRub)
	assert.Nil(t, err)

	_, err = NewTransaction("t-1", []Posting{postingUsd, postingRub}, someTime, "key-1", "")
	assert.ErrorIs(t, err, ErrDifferentCurrency)
}

func TestTransaction_Posted(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currencyUsd, err := NewCurrency("USD")
	assert.Nil(t, err)

	money := NewMoney(someAmount, currencyUsd)
	moneyNegate := NewMoney(-someAmount, currencyUsd)

	postingUsd, err := NewPosting(someAccountID, money)
	assert.Nil(t, err)

	postingRub, err := NewPosting(someAccountID, moneyNegate)
	assert.Nil(t, err)

	tx, err := NewTransaction("t-1", []Posting{postingUsd, postingRub}, someTime, "key-1", "")
	assert.Nil(t, err)

	events := tx.PullEvents()
	assert.Equal(t, 1, len(events))
}

func TestTransaction_NotZeroSum(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currency, err := NewCurrency("USD")
	assert.Nil(t, err)

	p1, err := NewPosting(someAccountID, NewMoney(someAmount, currency))
	assert.Nil(t, err)

	p2, err := NewPosting(someAccountID, NewMoney(someAmount, currency))
	assert.Nil(t, err)

	_, err = NewTransaction(someTransactionID, []Posting{p1, p2}, someTime, "key-1", "")
	assert.ErrorIs(t, err, ErrNotZeroSum)
}

func TestTransaction_DefensiveCopyOnInput(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currency, err := NewCurrency("USD")
	assert.Nil(t, err)

	p1, err := NewPosting(someAccountID, NewMoney(someAmount, currency))
	assert.Nil(t, err)

	p2, err := NewPosting(someAccountID, NewMoney(-someAmount, currency))
	assert.Nil(t, err)

	postings := []Posting{p1, p2}

	tx, err := NewTransaction(someTransactionID, postings, someTime, "key-1", "")
	assert.Nil(t, err)

	// мутация исходного слайса не должна протечь внутрь транзакции
	postings[0] = p2

	got := tx.Postings()
	assert.Equal(t, p1, got[0])
}

func TestTransaction_DefensiveCopyOnRead(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currency, err := NewCurrency("USD")
	assert.Nil(t, err)

	p1, err := NewPosting(someAccountID, NewMoney(someAmount, currency))
	assert.Nil(t, err)

	p2, err := NewPosting(someAccountID, NewMoney(-someAmount, currency))
	assert.Nil(t, err)

	tx, err := NewTransaction(someTransactionID, []Posting{p1, p2}, someTime, "key-1", "")
	assert.Nil(t, err)

	// мутация возвращённой копии не должна менять внутреннее состояние
	got := tx.Postings()
	got[0] = p2

	again := tx.Postings()
	assert.Equal(t, p1, again[0])
}

func TestTransaction_Accessors(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	currency, err := NewCurrency("USD")
	assert.Nil(t, err)

	p1, err := NewPosting(someAccountID, NewMoney(someAmount, currency))
	assert.Nil(t, err)

	p2, err := NewPosting(someAccountID, NewMoney(-someAmount, currency))
	assert.Nil(t, err)

	tx, err := NewTransaction(someTransactionID, []Posting{p1, p2}, someTime, "key-1", "")
	assert.Nil(t, err)

	assert.Equal(t, TransactionID(someTransactionID), tx.ID())
	assert.Equal(t, someTime, tx.CreatedAt())
	assert.Equal(t, IdempotencyKey("key-1"), tx.Key())
	assert.Equal(t, 2, len(tx.Postings()))
}
