package app

import (
	"FinFlow/internal/ledger/domain"
	"FinFlow/internal/ledger/infra/memory"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type fakeIDGen struct {
	n int
}

func (g *fakeIDGen) Generate() domain.TransactionID {
	g.n++
	return domain.TransactionID(fmt.Sprintf("tx-%d", g.n))
}

func TestTransfer_HappyPath(t *testing.T) {

	someTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	c, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenUserAccount("acc-1", "u-1", c, someTime)
	assert.Nil(t, err)

	to, err := domain.OpenUserAccount("acc-2", "u-2", c, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)
	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	ts := TransferCommand{
		From:        "acc-1",
		To:          "acc-2",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	tx, events, err := service.Transfer(t.Context(), ts)
	assert.Nil(t, err)

	assert.NotNil(t, tx)

	gotFrom, _ := accounts.FindByID(t.Context(), "acc-1")
	gotTo, _ := accounts.FindByID(t.Context(), "acc-2")

	assert.Equal(t, int64(-1000), gotFrom.Balance().Amount())
	assert.Equal(t, int64(1000), gotTo.Balance().Amount())

	assert.Equal(t, len(events), 1)
}

func TestTransfer_IdempotencyKey(t *testing.T) {

	someTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	c, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenUserAccount("acc-1", "u-1", c, someTime)
	assert.Nil(t, err)

	to, err := domain.OpenUserAccount("acc-2", "u-2", c, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)
	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	ts := TransferCommand{
		From:        "acc-1",
		To:          "acc-2",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	tx, events, err := service.Transfer(t.Context(), ts)
	assert.Nil(t, err)
	assert.NotNil(t, tx)

	gotFrom, _ := accounts.FindByID(t.Context(), "acc-1")
	gotTo, _ := accounts.FindByID(t.Context(), "acc-2")

	assert.Equal(t, int64(-1000), gotFrom.Balance().Amount())
	assert.Equal(t, int64(1000), gotTo.Balance().Amount())

	ts = TransferCommand{
		From:        "acc-1",
		To:          "acc-2",
		Amount:      1500,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	tx, events, err = service.Transfer(t.Context(), ts)
	assert.Nil(t, err)
	assert.NotNil(t, tx)

	gotFrom, _ = accounts.FindByID(t.Context(), "acc-1")
	gotTo, _ = accounts.FindByID(t.Context(), "acc-2")

	assert.Equal(t, int64(-1000), gotFrom.Balance().Amount())
	assert.Equal(t, int64(1000), gotTo.Balance().Amount())

	assert.Equal(t, 0, len(events))
}

func TestTransfer_Overdraft(t *testing.T) {

	someTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	c, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenUserAccount("acc-1", "u-1", c, someTime)
	assert.Nil(t, err)

	to, err := domain.OpenUserAccount("acc-2", "u-2", c, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)
	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	ts := TransferCommand{
		From:        "acc-1",
		To:          "acc-2",
		Amount:      6000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, events, err := service.Transfer(t.Context(), ts)
	assert.ErrorIs(t, err, domain.ErrOverdraftLimit)

	assert.Equal(t, 0, len(events))
}

func TestTransfer_InvalidAmount(t *testing.T) {
	someTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	c, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenUserAccount("acc-1", "u-1", c, someTime)
	assert.Nil(t, err)

	to, err := domain.OpenUserAccount("acc-2", "u-2", c, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)
	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	ts := TransferCommand{
		From:        "acc-1",
		To:          "acc-2",
		Amount:      -100,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, events, err := service.Transfer(t.Context(), ts)
	assert.ErrorIs(t, err, domain.ErrInvalidAmount)

	assert.Equal(t, 0, len(events))
}

func TestTransfer_SameAccount(t *testing.T) {

	someTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	c, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenUserAccount("acc-1", "u-1", c, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	ts := TransferCommand{
		From:        "acc-1",
		To:          "acc-1",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, events, err := service.Transfer(t.Context(), ts)
	assert.ErrorIs(t, err, domain.ErrSameAccount)
	assert.Equal(t, 0, len(events))
}

func TestTransfer_MismatchCurrency(t *testing.T) {

	someTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	usd, err := domain.NewCurrency("USD")
	assert.Nil(t, err)
	rub, err := domain.NewCurrency("RUB")
	assert.Nil(t, err)

	from, err := domain.OpenUserAccount("acc-1", "u-1", usd, someTime)
	assert.Nil(t, err)

	to, err := domain.OpenUserAccount("acc-2", "u-2", rub, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)
	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	ts := TransferCommand{
		From:        "acc-1",
		To:          "acc-2",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, events, err := service.Transfer(t.Context(), ts)
	assert.ErrorIs(t, err, domain.ErrMismatchCurrency)
	assert.Equal(t, 0, len(events))
}

func TestTransfer_NotFound(t *testing.T) {

	someTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	c, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	to, err := domain.OpenUserAccount("acc-2", "u-2", c, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	ts := TransferCommand{
		From:        "acc-1",
		To:          "acc-2",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, events, err := service.Transfer(t.Context(), ts)
	assert.ErrorIs(t, err, domain.ErrNotFound)
	assert.Equal(t, 0, len(events))
}
