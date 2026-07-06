package app

import (
	"FinFlow/internal/ledger/domain"
	"FinFlow/internal/ledger/infra/memory"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDeposit_HappyPath(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	c, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenSystemAccount("acc-1", c, someTime)
	assert.Nil(t, err)

	to, err := domain.OpenUserAccount("acc-2", "u-1", c, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)

	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	dc := DepositCommand{
		To:          "acc-2",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	tx, events, err := service.Deposit(t.Context(), dc)
	assert.Nil(t, err)
	assert.NotNil(t, tx)
	assert.NotNil(t, events)

	gotFrom, err := accounts.FindSystemAccount(t.Context(), from.Currency())
	assert.Nil(t, err)

	gotTo, err := accounts.FindByID(t.Context(), to.ID())
	assert.Nil(t, err)

	assert.Equal(t, int64(-1000), gotFrom.Balance().Amount())
	assert.Equal(t, int64(1000), gotTo.Balance().Amount())

	assert.Equal(t, 1, len(events))
}

func TestDeposit_IdempotencyKey(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	c, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenSystemAccount("acc-1", c, someTime)
	assert.Nil(t, err)

	to, err := domain.OpenUserAccount("acc-2", "u-1", c, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)

	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	dc := DepositCommand{
		To:          "acc-2",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	tx, events, err := service.Deposit(t.Context(), dc)
	assert.Nil(t, err)
	assert.NotNil(t, tx)
	assert.NotNil(t, events)

	dc = DepositCommand{
		To:          "acc-2",
		Amount:      1500,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	tx, events, err = service.Deposit(t.Context(), dc)
	assert.Nil(t, err)
	assert.Nil(t, events)
	assert.NotNil(t, tx)

	gotFrom, err := accounts.FindSystemAccount(t.Context(), from.Currency())
	assert.Nil(t, err)

	gotTo, err := accounts.FindByID(t.Context(), to.ID())
	assert.Nil(t, err)

	assert.Equal(t, int64(-1000), gotFrom.Balance().Amount())
	assert.Equal(t, int64(1000), gotTo.Balance().Amount())

	assert.Equal(t, 0, len(events))
}

func TestDeposit_NegativeAmount(t *testing.T) {
	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	c, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenSystemAccount("acc-1", c, someTime)
	assert.Nil(t, err)

	to, err := domain.OpenUserAccount("acc-2", "u-1", c, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)

	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	dc := DepositCommand{
		To:          "acc-2",
		Amount:      -1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, _, err = service.Deposit(t.Context(), dc)
	assert.ErrorIs(t, err, domain.ErrInvalidAmount)
}

func TestDeposit_NotUserAccount(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	c, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	to, err := domain.OpenSystemAccount("acc-1", c, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	dc := DepositCommand{
		To:          "acc-1",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, _, err = service.Deposit(t.Context(), dc)
	assert.ErrorIs(t, err, domain.ErrNotUserAccount)
}

func TestDeposit_ToNotFound(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	dc := DepositCommand{
		To:          "acc-1",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, _, err := service.Deposit(t.Context(), dc)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestDeposit_SystemAccountNotFound(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	c, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	to, err := domain.OpenUserAccount("acc-1", "u-1", c, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	dc := DepositCommand{
		To:          "acc-1",
		Amount:      1000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	_, _, err = service.Deposit(t.Context(), dc)
	assert.ErrorIs(t, err, domain.ErrSystemAccountNotFound)
}

func TestDeposit_LargeAmount(t *testing.T) {

	someTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	c, err := domain.NewCurrency("USD")
	assert.Nil(t, err)

	from, err := domain.OpenSystemAccount("acc-1", c, someTime)
	assert.Nil(t, err)

	to, err := domain.OpenUserAccount("acc-2", "u-1", c, someTime)
	assert.Nil(t, err)

	accounts := memory.NewAccountRepository()
	txns := memory.NewTransactionRepository()

	err = accounts.Save(t.Context(), from)
	assert.Nil(t, err)
	err = accounts.Save(t.Context(), to)
	assert.Nil(t, err)

	service := NewLedgerService(accounts, txns, &fakeIDGen{})

	dc := DepositCommand{
		To:          "acc-2",
		Amount:      1_000_000,
		Key:         "key-1",
		At:          someTime,
		Description: "",
	}

	tx, events, err := service.Deposit(t.Context(), dc)
	assert.Nil(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, 1, len(events))

	gotFrom, err := accounts.FindSystemAccount(t.Context(), c)
	assert.Nil(t, err)

	gotTo, err := accounts.FindByID(t.Context(), "acc-2")
	assert.Nil(t, err)

	assert.Equal(t, int64(-1_000_000), gotFrom.Balance().Amount())
	assert.Equal(t, int64(1_000_000), gotTo.Balance().Amount())
}
